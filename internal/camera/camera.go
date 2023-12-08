package camera

import (
	math "github.com/chewxy/math32"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	position mgl32.Vec3
	front    mgl32.Vec3
	up       mgl32.Vec3

	direction mgl32.Vec3
	movement  mgl32.Vec3

	projection mgl32.Mat4
	view       mgl32.Mat4

	yaw         float32
	pitch       float32
	speed       float32
	sensitivity float32
	sprint      bool
}

func NewCamera(pos mgl32.Vec3, windowWidth, windowHeight int) *Camera {
	c := &Camera{
		position:    pos,
		front:       mgl32.Vec3{0.0, 0.0, 1.0},
		up:          mgl32.Vec3{0.0, 1.0, 0.0},
		movement:    mgl32.Vec3{0.0, 0.0, 0.0},
		speed:       3.0,
		sensitivity: 0.075,
		yaw:         -90.0,
		pitch:       -10.0,
	}

	c.projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 1000.0)

	return c
}

func (c *Camera) Update(dt float32) {
	// Look direction
	c.direction[0] = math.Cos(mgl32.DegToRad(c.yaw)) * math.Cos(mgl32.DegToRad(c.pitch))
	c.direction[1] = math.Sin(mgl32.DegToRad(c.pitch))
	c.direction[2] = math.Sin(mgl32.DegToRad(c.yaw)) * math.Cos(mgl32.DegToRad(c.pitch))
	c.front = c.direction.Normalize()

	// Sprint
	speed := c.speed * dt
	if c.sprint {
		speed = speed * 3
	}

	// forward/backward
	if c.movement.Z() > 0.5 {
		c.position = c.position.Add(c.front.Mul(speed))
	} else if c.movement.Z() < -0.5 {
		c.position = c.position.Sub(c.front.Mul(speed))
	}

	// left/right
	if c.movement.X() > 0.5 {
		c.position = c.position.Sub(c.front.Cross(c.up).Normalize().Mul(speed))
	} else if c.movement.X() < -0.5 {
		c.position = c.position.Add(c.front.Cross(c.up).Normalize().Mul(speed))
	}

	// up/down
	if c.movement.Y() > 0.5 {
		c.position = c.position.Add(c.up.Mul(speed))
	} else if c.movement.Y() < -0.5 {
		c.position = c.position.Sub(c.up.Mul(speed))
	}

	c.view = mgl32.LookAtV(c.position, c.position.Add(c.front), c.up)
}

func (c *Camera) GetSensitivity() float32 {
	return c.sensitivity
}

func (c *Camera) SetSensitivity(sens float32) {
	c.sensitivity = sens
}

func (c *Camera) GetSpeed() float32 {
	return c.speed
}

func (c *Camera) SetSpeed(speed float32) {
	c.speed = speed
}

func (c *Camera) Sprint(sprint bool) {
	c.sprint = sprint
}

func (c *Camera) AddYaw(yaw float32) {
	c.yaw += yaw * c.sensitivity
}

func (c *Camera) AddPitch(pitch float32) {
	c.pitch += pitch * c.sensitivity
	c.pitch = mgl32.Clamp(c.pitch, -89.5, 89.5)
}

func (c *Camera) View() mgl32.Mat4 {
	return c.view
}

func (c *Camera) Projection() mgl32.Mat4 {
	return c.projection
}

func (c *Camera) MoveX(x float32) {
	c.movement[0] += x
}

func (c *Camera) MoveY(y float32) {
	c.movement[1] += y
}

func (c *Camera) MoveZ(z float32) {
	c.movement[2] += z
}
