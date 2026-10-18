package canvas3d

import "time"

var (
	startTime   = time.Now()
	lastFrame   = startTime
	deltaTime   float32
	elapsedTime float32
)

// Call once per frame.
func UpdateTime() {
	now := time.Now()

	deltaTime = float32(now.Sub(lastFrame).Seconds())
	elapsedTime = float32(now.Sub(startTime).Seconds())

	lastFrame = now
}

// Time since the previous UpdateTime call, in seconds.
func DeltaTime() float32 {
	return deltaTime
}

// Total time since the package initialized, in seconds.
func ElapsedTime() float32 {
	return elapsedTime
}
