package main

import (
	"fmt"
	"log"
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

//nolint:funlen,gocognit,gocyclo,maintidx // foo
func main() {
	runtime.LockOSThread()

	// Init GLFW/OpenGL
	window, clean := initGlfw()
	defer clean()
	initOpenGL()

	camera := camera.NewCamera(mgl32.Vec3{0.0, 2.0, 10.0}, windowWidth, windowHeight)

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

	thing := makeThing()

	cube := makeCube()
	grid := makeGrid()

	// --------------------------------------------------------------------------------------------

	// angle := 0.0
	previousTime := float32(glfw.GetTime())

	// Game loop
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Calculate time since last frame
		time := float32(glfw.GetTime())
		dt := time - previousTime
		previousTime = time

		// angle += elapsed
		// model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// Update resources
		camera.Update(dt)
		cube.Update(dt)
		thing.Update(dt)

		// Render resources
		cube.Render(camera)
		thing.Render2(camera)

		// Draw grid last for some reason? Why is this?
		grid.Render(camera)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func makeThing() *gameobject.GameObject {
	wd, _ := os.Getwd()
	thing, err := mesh.FromFile(wd + "/resources/meshes/cube.obj")
	if err != nil {
		log.Fatal("error making thing", err)
	}

	shader, err := shader.NewSolidShader()
	if err != nil {
		log.Fatal(err)
	}

	obj := gameobject.NewBuilder().
		Position(mgl32.Vec3{0.0, 3.0, 0.0}).
		// Scale(mgl32.Vec3{0.05, 0.05, 0.05}).
		Mesh(&thing).Shader(&shader).Build()

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
	gl.DebugMessageCallback(func(
		_ uint32, _ uint32, _ uint32, _ uint32, _ int32, message string, _ unsafe.Pointer,
	) {
		fmt.Printf("OPENGL MESSAGE: %v\n", message)
	}, nil)

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.DepthFunc(gl.LESS)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.2, 0.2, 0.2, 1.0)
}
