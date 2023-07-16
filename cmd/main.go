package main

import (
	"fmt"
	"log"
	"runtime"

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

			// Z to toggle between full, wireframe, or vertex render
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

	// --- CUBE -----------------------------------------------------------------------------------

	// vertexShaderFilename := filepath.Join(wd, "../resources/shaders/texture.vert")
	// fragmentShaderFilename := filepath.Join(wd, "../resources/shaders/texture.frag")

	// cubeMesh := mesh.MakeCube()

	// cubeShader, err := shader.NewBasicShader()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// squareTexture, err := texture.New("../resources/textures/square.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// cube := gameobject.GameObject{
	// 	Position: mgl32.Vec3{0.0, 0.0, 0.0},
	// 	Rotation: mgl32.Vec3{0.0, 0.0, 0.0},
	// 	Scale:    mgl32.Vec3{1.0, 1.0, 1.0},

	// 	Shader:  &cubeShader,
	// 	Texture: &squareTexture,
	// 	Mesh:    &cubeMesh,
	// }

	cube := makeCube()
	grid := makeGrid()

	// cubeShader, err := shader.New(vertexShaderFile, fragmentShaderFile)
	// if err != nil {
	// 	fmt.Printf("Error creating shader: %v\n", err)
	// 	panic(err)
	// }
	// shaders = append(shaders, cubeShader)

	// Setup
	// cubeShader.SetProjection(camera.Projection())
	// cubeShader.SetView(camera.View())

	// modelMat := mgl32.Ident4()
	// cubeShader.SetModel(modelMat)

	// program := cubeShader.GetProgram()

	// textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	// gl.Uniform1i(textureUniform, 0)

	// gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// texture, err := texture.New("../resources/textures/square.png")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// var vao uint32
	// gl.GenVertexArrays(1, &vao)
	// gl.BindVertexArray(vao)

	// var vbo uint32
	// gl.GenBuffers(1, &vbo)
	// gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// gl.BufferData(gl.ARRAY_BUFFER, len(mesh.Cube)*4, gl.Ptr(mesh.Cube), gl.STATIC_DRAW)

	// vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	// gl.EnableVertexAttribArray(vertAttrib)
	// gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)

	// cubeTexCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	// gl.EnableVertexAttribArray(cubeTexCoordAttrib)
	// gl.VertexAttribPointerWithOffset(cubeTexCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)

	// --- TRIANGLE -------------------------------------------------------------------------------

	// triangleShader, err := shader.New(vertexShaderFile, fragmentShaderFile)
	// if err != nil {
	// 	fmt.Printf("Error creating shader: %v\n", err)
	// 	panic(err)
	// }
	// shaders = append(shaders, triangleShader)

	// Setup
	// 	triangleShader.SetProjection(camera.Projection())
	//
	// 	triangleShader.SetView(camera.View())
	//
	// 	triangleModel := mgl32.Ident4()
	// 	triangleModel = triangleModel.Mul4(mgl32.Translate3D(3.0, 3.0, 0.0))
	// 	// triangleModel.Set(1, 3, 3.0)
	// 	triangleShader.SetModel(triangleModel)
	//
	// 	triangleProgram := triangleShader.GetProgram()
	//
	// 	triangletextureUniform := gl.GetUniformLocation(triangleProgram, gl.Str("tex\x00"))
	// 	gl.Uniform1i(triangletextureUniform, 0)
	//
	// 	gl.BindFragDataLocation(triangleProgram, 0, gl.Str("outputColor\x00"))
	//
	// 	var triangleVao uint32
	// 	gl.GenVertexArrays(1, &triangleVao)
	// 	gl.BindVertexArray(triangleVao)
	//
	// 	var triangleVbo uint32
	// 	gl.GenBuffers(1, &triangleVbo)
	// 	gl.BindBuffer(gl.ARRAY_BUFFER, triangleVbo)
	// 	gl.BufferData(gl.ARRAY_BUFFER, len(model.Triangle)*4, gl.Ptr(model.Triangle), gl.STATIC_DRAW)
	//
	// 	triangleVertAttrib := uint32(gl.GetAttribLocation(triangleProgram, gl.Str("vert\x00")))
	// 	gl.EnableVertexAttribArray(triangleVertAttrib)
	// 	gl.VertexAttribPointerWithOffset(triangleVertAttrib, 3, gl.FLOAT, false, 5*4, 0)
	//
	// 	triangleTexCoordAttrib := uint32(gl.GetAttribLocation(triangleProgram, gl.Str("vertTexCoord\x00")))
	// 	gl.EnableVertexAttribArray(triangleTexCoordAttrib)
	// 	gl.VertexAttribPointerWithOffset(triangleTexCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)

	// --- GRID -----------------------------------------------------------------------------------

	// 	gridShader, err := shader.NewGridShader()
	// 	if err != nil {
	// 		fmt.Printf("Error creating grid shader: %v\n", err)
	// 		panic(err)
	// 	}
	// 	shaders = append(shaders, &gridShader)
	//
	// 	gridShader.SetProjection(camera.Projection())
	// 	gridShader.SetView(camera.View())
	//
	// 	gridModel := mgl32.Ident4()
	// 	// gridModel = gridModel.Mul4(mgl32.Scale3D(3, 3, 1))
	// 	// gridModel = gridModel.Mul4(mgl32.HomogRotate3DX(mgl32.DegToRad(90)))
	// 	gridShader.SetModel(gridModel)
	//
	// 	gridProgram := gridShader.GetProgram()
	//
	// 	// gl.BindFragDataLocation(gridProgram, 0, gl.Str("outputColor\x00"))
	//
	// 	var gridVao uint32
	// 	gl.GenVertexArrays(1, &gridVao)
	// 	gl.BindVertexArray(gridVao)
	//
	// 	var gridVbo uint32
	// 	gl.GenBuffers(1, &gridVbo)
	// 	gl.BindBuffer(gl.ARRAY_BUFFER, gridVbo)
	// 	gl.BufferData(gl.ARRAY_BUFFER, len(mesh.Plane)*4, gl.Ptr(mesh.Plane), gl.STATIC_DRAW)
	//
	// 	gridVertAttrib := uint32(gl.GetAttribLocation(gridProgram, gl.Str("vert\x00")))
	// 	gl.EnableVertexAttribArray(gridVertAttrib)
	// 	gl.VertexAttribPointerWithOffset(gridVertAttrib, 3, gl.FLOAT, false, 5*4, 0)
	//
	// 	gridTexCoordAttrib := uint32(gl.GetAttribLocation(gridProgram, gl.Str("uv\x00")))
	// 	gl.EnableVertexAttribArray(gridTexCoordAttrib)
	// 	gl.VertexAttribPointerWithOffset(gridTexCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)

	// --------------------------------------------------------------------------------------------

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.DepthFunc(gl.LESS)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.2, 0.2, 0.2, 1.0)

	// angle := 0.0
	previousTime := float32(glfw.GetTime())

	// Game loop
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := float32(glfw.GetTime())
		dt := time - previousTime
		previousTime = time

		// angle += elapsed
		// model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		camera.Update(dt)

		// Update
		cube.Update(dt)

		// Render cube
		cube.Render(camera)
		grid.Render(camera)

		// gl.UseProgram(cubeShader.GetProgram())
		// cubeShader.SetModel(modelMat)
		// cubeShader.SetView(camera.View())

		// gl.BindVertexArray(cube.Mesh.Vao)
		// gl.BindBuffer(gl.ARRAY_BUFFER, cube.Mesh.Vbo)

		// gl.ActiveTexture(gl.TEXTURE0)
		// gl.BindTexture(gl.TEXTURE_2D, squareTexture.Handle)

		// gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)

		// gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		// done with cube

		// render triangle
		// 		gl.UseProgram(triangleShader.GetProgram())
		// 		triangleShader.SetModel(triangleModel)
		// 		triangleShader.SetView(camera.View())
		//
		// 		gl.BindVertexArray(triangleVao)
		// 		gl.BindBuffer(gl.ARRAY_BUFFER, triangleVbo)
		//
		// 		gl.ActiveTexture(gl.TEXTURE0)
		// 		gl.BindTexture(gl.TEXTURE_2D, texture)
		//
		// 		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		//
		// 		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		// done with triangle

		// render grid
		// 		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		// 		gl.UseProgram(gridShader.GetProgram())
		// 		gridShader.SetModel(gridModel)
		// 		gridShader.SetView(camera.View())
		//
		// 		gl.BindVertexArray(gridVao)
		// 		gl.BindBuffer(gl.ARRAY_BUFFER, gridVbo)
		//
		// 		gl.DrawArrays(gl.TRIANGLES, 0, 3*2)
		//
		// 		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		// 		gl.PolygonMode(gl.FRONT_AND_BACK, drawMode)
		// done with grid

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func makeCube() gameobject.GameObject {
	cubeMesh := mesh.MakeCube()

	cubeShader, err := shader.NewBasicShader()
	if err != nil {
		log.Fatal(err)
	}

	squareTexture, err := texture.New("../resources/textures/square.png")
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
