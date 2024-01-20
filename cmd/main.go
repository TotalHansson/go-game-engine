package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"unsafe"

	"game-engine/rts/internal/camera"
	"game-engine/rts/internal/gameobject"
	"game-engine/rts/internal/mesh"
	"game-engine/rts/internal/shader"
	"game-engine/rts/internal/texture"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 1280
	windowHeight = 720
)

var (
	shaders         = []*shader.Shader{}
	drawMode uint32 = gl.FILL
)

type State interface {
	OnEnter()
	OnUpdate(dt float32)
	OnLeave()
}

type EmptyState struct{}

func (s *EmptyState) OnEnter()            {}
func (s *EmptyState) OnUpdate(dt float32) {}
func (s *EmptyState) OnLeave()            {}

type BasicState struct {
	name        string
	timeInState float32
}

func (s *BasicState) OnEnter() {
	fmt.Printf("Entering state %q\n", s.name)
}
func (s *BasicState) OnUpdate(dt float32) {
	s.timeInState += dt
	fmt.Printf("In state %s for %v\n", s.name, s.timeInState)
}
func (s *BasicState) OnLeave() {
	fmt.Printf("Leaving state %q\n", s.name)
}
func NewBasicState(name string) BasicState {
	return BasicState{
		name: name,
	}
}

type Worker struct {
	fsm              FSM
	idleState        WorkerIdleState
	choppingState    WorkerChoppingState
	treePickingState WorkerTreePickingState
	walkState        WorkerWalkState

	trees         *[]*gameobject.SolidGameObject
	gameObject    *gameobject.SolidGameObject
	currentTarget *gameobject.SolidGameObject
}

func (w *Worker) Idle() {
	w.fsm.ChangeState(&w.idleState)
}
func (w *Worker) Walk() {
	w.fsm.ChangeState(&w.walkState)
}
func (w *Worker) Chop() {
	w.fsm.ChangeState(&w.choppingState)
}
func (w *Worker) PickTree() {
	w.fsm.ChangeState(&w.treePickingState)
}

type WorkerIdleState struct {
	worker             *Worker
	timeSinceLastPrint float32
}

func (s *WorkerIdleState) OnUpdate(dt float32) {
	s.timeSinceLastPrint += dt
	if s.timeSinceLastPrint > 1.0 {
		s.timeSinceLastPrint = 0.0
		fmt.Printf("Worker is chillin\n")
	}
}
func (s *WorkerIdleState) OnEnter() {
	s.timeSinceLastPrint = 0.0
	fmt.Printf("Entering idle state\n")
}
func (s *WorkerIdleState) OnLeave() { fmt.Printf("Leaving idle state\n") }

var treeChoppingDuration float32 = 10.0

type WorkerChoppingState struct {
	worker             *Worker
	timeSinceLastPrint float32
	timeSpentChopping  float32
}

func (s *WorkerChoppingState) OnUpdate(dt float32) {
	s.timeSinceLastPrint += dt
	if s.timeSinceLastPrint > 1.0 {
		s.timeSinceLastPrint = 0
		fmt.Printf("Chopping tree: %.1f/%.1f\n", s.timeSpentChopping, treeChoppingDuration)
	}

	s.timeSpentChopping += dt
	if s.timeSpentChopping >= treeChoppingDuration {
		fmt.Printf("Timber!\n")
		*s.worker.trees, _ = chopTree(s.worker.currentTarget, *s.worker.trees)
		s.worker.PickTree()
	}
}
func (s *WorkerChoppingState) OnEnter() {
	s.timeSpentChopping = 0
	s.timeSinceLastPrint = 0
	fmt.Printf("Entering chopping state\n")
}
func (s *WorkerChoppingState) OnLeave() { fmt.Printf("Leaving chopping state\n") }

type WorkerTreePickingState struct {
	worker *Worker
}

func (s *WorkerTreePickingState) OnUpdate(dt float32) {
	targetIndex, exists := findClosestTree(*s.worker.trees, s.worker.gameObject)
	if !exists {
		fmt.Printf("No more trees\n")
		s.worker.Idle()
		return
	}
	trees := *s.worker.trees
	s.worker.currentTarget = trees[targetIndex]
	s.worker.Walk()
}
func (s *WorkerTreePickingState) OnEnter() { fmt.Printf("Entering tree picking state\n") }
func (s *WorkerTreePickingState) OnLeave() { fmt.Printf("Leaving tree picking state\n") }

type WorkerWalkState struct {
	worker             *Worker
	timeSinceLastPrint float32
}

func (s *WorkerWalkState) OnUpdate(dt float32) {
	target := s.worker.currentTarget
	worker := s.worker.gameObject
	moveDir := target.Position.Sub(worker.Position)
	moveDir[1] = 0.0 // Don't touch the height
	moveDir = moveDir.Normalize()
	worker.Position = worker.Position.Add(moveDir.Mul(dt))
	dist := worker.Position.Sub(target.Position)
	dist[1] = 0.0
	remainingDistance := dist.Len()
	if remainingDistance < 0.05 {
		s.worker.Chop()
		return
	}

	s.timeSinceLastPrint += dt
	if s.timeSinceLastPrint > 1.0 {
		s.timeSinceLastPrint = 0.0
		fmt.Printf("Walking towards tree %v, remaining distance: %v\n", s.worker.currentTarget, remainingDistance)
	}
}
func (s *WorkerWalkState) OnEnter() { fmt.Printf("Entering walking state\n") }
func (s *WorkerWalkState) OnLeave() { fmt.Printf("Leaving walking state\n") }

type FSM struct {
	currentState State
}

func (fsm *FSM) Run(dt float32) {
	fsm.currentState.OnUpdate(dt)
}
func (fsm *FSM) ChangeState(newState State) {
	fsm.currentState.OnLeave()
	fsm.currentState = newState
	fsm.currentState.OnEnter()
}
func NewFSM() *FSM {
	return &FSM{currentState: &EmptyState{}}
}

//nolint:funlen,gocognit,gocyclo,maintidx // foo
func main() {
	runtime.LockOSThread()

	// Init GLFW/OpenGL
	window, clean := initGlfw()
	defer clean()
	initOpenGL()

	camera := camera.NewCamera(mgl32.Vec3{4.0, 4.0, 10.0}, windowWidth, windowHeight)

	keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		switch key {
		case glfw.KeyEscape:
			window.SetShouldClose(true)

		case glfw.KeyW:
			if action == glfw.Press {
				camera.MoveZ(1)
			} else if action == glfw.Release {
				camera.MoveZ(-1)
			}

		case glfw.KeyS:
			if action == glfw.Press {
				camera.MoveZ(-1)
			} else if action == glfw.Release {
				camera.MoveZ(1)
			}

		case glfw.KeyA:
			if action == glfw.Press {
				camera.MoveX(1)
			} else if action == glfw.Release {
				camera.MoveX(-1)
			}

		case glfw.KeyD:
			if action == glfw.Press {
				camera.MoveX(-1)
			} else if action == glfw.Release {
				camera.MoveX(1)
			}

		case glfw.KeyE:
			if action == glfw.Press {
				camera.MoveY(1)
			} else if action == glfw.Release {
				camera.MoveY(-1)
			}

		case glfw.KeyQ:
			if action == glfw.Press {
				camera.MoveY(-1)
			} else if action == glfw.Release {
				camera.MoveY(1)
			}

		case glfw.KeyLeftShift, glfw.KeyRightShift:
			if action == glfw.Press {
				camera.Sprint(true)
			} else if action == glfw.Release {
				camera.Sprint(false)
			}

		case glfw.KeyR:
			if action == glfw.Press {
				reloadShaders()
			}

			// Z to toggle between full or wireframe render
		case glfw.KeyZ:
			if action == glfw.Press {
				if drawMode == gl.FILL {
					drawMode = gl.LINE
				} else if drawMode == gl.LINE {
					drawMode = gl.FILL
				}
				gl.PolygonMode(gl.FRONT_AND_BACK, drawMode)
			}
		}
	}

	var previousMouseX, previousMouseY float32 = windowWidth / 2.0, windowHeight / 2.0

	var xOffset, yOffset float32 = windowWidth / 2.0, windowHeight / 2.0

	firstMouse := true

	mouseCallback := func(w *glfw.Window, xpos float64, ypos float64) {
		if firstMouse {
			previousMouseX, previousMouseY = float32(xpos), float32(ypos)
			firstMouse = false
		}

		xOffset = float32(xpos) - previousMouseX
		yOffset = previousMouseY - float32(ypos)
		previousMouseX, previousMouseY = float32(xpos), float32(ypos)

		camera.AddYaw(xOffset)
		camera.AddPitch(yOffset)
	}

	window.SetKeyCallback(keyCallback)
	window.SetCursorPosCallback(mouseCallback)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	setGlobalGLState()
	wd, _ := os.Getwd()

	///////////// XYZ Gizmo //////////////
	xyzShader, err := shader.NewXYZGizmoShader()
	if err != nil {
		log.Fatal(err)
	}
	xyzMesh, err := mesh.FromFile(wd + "/resources/meshes/xyz-gizmo.obj")
	if err != nil {
		log.Fatal("error making thing", err)
	}
	xyz := &gameobject.XYZGizmo{
		Position: mgl32.Vec3{0.0, 0.0, -1.0},
		Scale:    mgl32.Vec3{1.0, 1.0, 1.0},
		Mesh:     &xyzMesh,
		Shader:   &xyzShader,
	}

	//////////// Light ///////////////

	lampPos := mgl32.Vec3{10.0, 10.0, -10.0}
	lampColor := mgl32.Vec3{1.0, 1.0, 1.0}

	s, err := shader.NewSolidShader(mgl32.Vec3{0.2, 0.4, 0.2})
	if err != nil {
		log.Fatal(err)
	}
	s.SetLightPos(lampPos)
	s.SetLightColor(lampColor)

	//////////// Ground //////////////

	bevelCube, err := mesh.FromFile(wd + "/resources/meshes/bevel-cube2.obj")
	if err != nil {
		log.Fatal("error making thing", err)
	}

	sizeX, sizeZ := 10, 10
	land := make([][]*gameobject.SolidGameObject, sizeX)
	for x := 0; x < sizeX; x++ {
		land[x] = make([]*gameobject.SolidGameObject, sizeZ)
		for z := 0; z < sizeZ; z++ {
			land[x][z] = &gameobject.SolidGameObject{
				Position: mgl32.Vec3{float32(x) * 2.0, 2.0, float32(z) * -2.0},
				Scale:    mgl32.Vec3{1.0, 0.5, 1.0},
				Mesh:     &bevelCube,
				Shader:   &s,
			}
		}
	}

	/////////// Trees ///////////////

	treeShader, err := shader.NewSolidShader(mgl32.Vec3{0.0, 1.0, 0.0})
	if err != nil {
		log.Fatal(err)
	}
	s.SetLightPos(lampPos)
	s.SetLightColor(lampColor)

	treeMesh, err := mesh.FromFile(wd + "/resources/meshes/cylinder.obj")
	if err != nil {
		log.Fatal("error loading tree mesh", err)
	}
	totalNrTrees := 500
	trees := make([]*gameobject.SolidGameObject, totalNrTrees)
	for i := 0; i < totalNrTrees; i++ {
		x := rand.Float32()*9.5 - 0.25
		z := rand.Float32()*9.5 - 0.25
		trees[i] = &gameobject.SolidGameObject{
			Position: mgl32.Vec3{x * 2.0, 3.0, z * -2.0},
			Scale:    mgl32.Vec3{0.1, 0.5, 0.1},
			Mesh:     &treeMesh,
			Shader:   &treeShader,
		}
	}

	/////////// Worker ///////////////

	workerShader, err := shader.NewSolidShader(mgl32.Vec3{0.75, 0.75, 0.75})
	if err != nil {
		log.Fatal(err)
	}
	workerShader.SetLightPos(lampPos)
	workerShader.SetLightColor(lampColor)
	workerMesh, err := mesh.FromFile(wd + "/resources/meshes/cube.obj")
	if err != nil {
		log.Fatal("error loading worker mesh", err)
	}
	worker := gameobject.SolidGameObject{
		Position: mgl32.Vec3{0.0, 2.5, 0.0},
		Scale:    mgl32.Vec3{0.2, 0.2, 0.2},
		Mesh:     &workerMesh,
		Shader:   &workerShader,
	}
	theAlmightyWorkerMan := &Worker{}
	theAlmightyWorkerMan.fsm = *NewFSM()
	theAlmightyWorkerMan.idleState = WorkerIdleState{worker: theAlmightyWorkerMan}
	theAlmightyWorkerMan.choppingState = WorkerChoppingState{worker: theAlmightyWorkerMan}
	theAlmightyWorkerMan.treePickingState = WorkerTreePickingState{worker: theAlmightyWorkerMan}
	theAlmightyWorkerMan.walkState = WorkerWalkState{worker: theAlmightyWorkerMan}

	theAlmightyWorkerMan.trees = &trees
	theAlmightyWorkerMan.gameObject = &worker
	theAlmightyWorkerMan.PickTree()
	// theAlmightyWorkerMan.currentTarget = trees[0]
	// closestTree, _ := findClosestTree(trees, &worker)

	//////////////////////////

	cube := makeCube()
	grid := makeGrid()

	// --------------------------------------------------------------------------------------------

	previousTime := float32(glfw.GetTime())

	// Game loop
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Viewport(0, 0, windowWidth, windowHeight)

		// Calculate time since last frame
		time := float32(glfw.GetTime())
		dt := time - previousTime
		previousTime = time

		// Update resources
		camera.Update(dt)
		cube.Update(dt)
		for x := 0; x < sizeX; x++ {
			for z := 0; z < sizeZ; z++ {
				land[x][z].Update(dt)
				land[x][z].Render(camera)
			}
		}
		for i := 0; i < len(trees); i++ {
			trees[i].Update(dt)
			trees[i].Render(camera)
		}

		//////// worker //////////
		// target := trees[closestTree]
		// moveDir := target.Position.Sub(worker.Position)
		// moveDir[1] = 0.0 // Don't touch the height
		// moveDir = moveDir.Normalize()
		// worker.Position = worker.Position.Add(moveDir.Mul(dt))
		// dist := worker.Position.Sub(target.Position)
		// dist[1] = 0.0
		// d := dist.Len()
		// if d < 0.05 {
		// 	trees, totalNrTrees = chopTree(closestTree, trees)
		// 	closestTree, _ = findClosestTree(trees, &worker)
		// }
		theAlmightyWorkerMan.fsm.Run(dt)

		worker.Update(dt)
		worker.Render(camera)
		//////////////////////////

		// Render resources
		cube.Render(camera)

		// Draw grid last for some reason? Why is this?
		grid.Render(camera)

		/////

		gl.Clear(gl.DEPTH_BUFFER_BIT)
		xyz.Update(dt)
		xyz.Render(camera)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func chopTree(tree *gameobject.SolidGameObject, trees []*gameobject.SolidGameObject) (remainingTrees []*gameobject.SolidGameObject, totalNrTrees int) {
	nrTrees := len(trees)
	index := 0
	for i := 0; i < nrTrees; i++ {
		if trees[i] == tree {
			index = i
			break
		}
	}
	trees[index], trees[nrTrees-1] = trees[nrTrees-1], trees[index]
	trees = trees[:nrTrees-1]
	return trees, nrTrees - 1
}

func findClosestTree(trees []*gameobject.SolidGameObject, worker *gameobject.SolidGameObject) (treeIndex int, targetExists bool) {
	nrTrees := len(trees)
	if nrTrees == 0 {
		return 0, false
	}
	closestTree := 0
	closestSqLen := trees[closestTree].Position.Sub(worker.Position).LenSqr()
	for i := 0; i < nrTrees; i++ {
		sqLen := trees[i].Position.Sub(worker.Position).LenSqr()
		if sqLen < closestSqLen {
			closestTree = i
			closestSqLen = sqLen
		}
	}
	return closestTree, true
}

func makeThing() *gameobject.GameObject {
	wd, _ := os.Getwd()
	thing, err := mesh.FromFile(wd + "/resources/meshes/cube.obj")
	if err != nil {
		log.Fatal("error making thing", err)
	}

	shader, err := shader.NewGenericShader()
	if err != nil {
		log.Fatal(err)
	}

	// texture, err := texture.New(wd + "/resources/textures/square.png")
	// if err != nil {
	// 	log.Fatal("error loading texture", err)
	// }

	obj := gameobject.NewBuilder().
		Position(mgl32.Vec3{0.0, 3.0, 0.0}).
		// Scale(mgl32.Vec3{0.05, 0.05, 0.05}).
		Mesh(&thing).Shader(&shader). //Texture(&texture).
		Build()

	return obj
}

func makeCube() gameobject.GameObject {
	cubeMesh := mesh.MakeCube()
	// wd, _ := os.Getwd()
	// cubeMesh, err := mesh.FromFile(wd + "/resources/mesh/cube.obj")

	cubeShader, err := shader.NewBasicShader()
	if err != nil {
		log.Fatal(err)
	}

	wd, _ := os.Getwd()
	squareTexture, err := texture.New(wd + "/resources/textures/square.png")
	if err != nil {
		log.Fatal(err)
	}

	cube := gameobject.GameObject{
		Position: mgl32.Vec3{0.0, 0.0, 0.0},
		Rotation: mgl32.Vec3{0.0, 0.0, 0.0},
		Scale:    mgl32.Vec3{1.0, 1.0, 1.0},

		Shader:  &cubeShader,
		Texture: &squareTexture,
		Mesh:    &cubeMesh,
	}

	return cube
}

func makeGrid() gameobject.GameObject {
	gridMesh := mesh.MakeGrid()

	gridShader, err := shader.NewGridShader()
	if err != nil {
		log.Fatal(err)
	}

	grid := gameobject.GameObject{
		Position: mgl32.Vec3{0.0, 0.0, 0.0},
		Rotation: mgl32.Vec3{0.0, 0.0, 0.0},
		Scale:    mgl32.Vec3{1.0, 1.0, 1.0},

		Shader: &gridShader,
		Mesh:   &gridMesh,
	}

	return grid
}

func reloadShaders() {
	fmt.Printf("Reloading shaders\n")

	for _, s := range shaders {
		if err := s.Reload(); err != nil {
			fmt.Printf("Error reloading shader: %v\n", err)
		}
	}
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() (*glfw.Window, func()) {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "My Window", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	cleanFunc := func() {
		glfw.Terminate()
	}

	return window, cleanFunc
}

// initOpenGL initializes OpenGL and returns an initialized program.
func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
}

func setGlobalGLState() {
	gl.Enable(gl.DEBUG_OUTPUT)

	// source uint32,
	// gltype uint32,
	// id uint32,
	// severity uint32,
	// length int32,
	// message *uint8,
	// userParam unsafe.Pointer
	gl.DebugMessageCallback(func(
		source uint32, glType uint32, id uint32, severity uint32, length int32, message string, _ unsafe.Pointer,
	) {
		fmt.Printf("OPENGL MESSAGE: %v, %v, %v, %v, %v - %v\n", source, glType, id, severity, length, message)
	}, nil)

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.DepthFunc(gl.LESS)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.2, 0.2, 0.2, 1.0)
	gl.Viewport(0, 0, windowWidth, windowHeight)
}
