package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game definitions

const (
	screenWidth  = 800
	screenHeight = 600
)

// Player definitions

type Player struct {
	Position rl.Vector2
	Speed    float32
	Rotation float32
}

// Enemy definitions

type Enemy struct {
	Position rl.Vector2
	Speed    float32
	Rotation float32
	Alive    bool
	ShootCounter float32
}

const (
	enemyRange       = 250.0
	enemyShootDelay  = 60.0 // roughly once a second if your FPS is 60
)

func InitializeEnemies() []Enemy {
    return []Enemy{
        {rl.NewVector2(500, 500), 2.0, 0.0, true, 0.0},
        {rl.NewVector2(3300, 500), 2.0, 0.0, true, 0.0},
        {rl.NewVector2(500, 2300), 2.0, 0.0, true, 0.0},
        {rl.NewVector2(600, 300), 2.0, 0.0, true, 0.0},
        {rl.NewVector2(300, 2500), 2.0, 0.0, true, 0.0},
        {rl.NewVector2(1300, 1300), 2.0, 0.0, true, 0.0},
        {rl.NewVector2(2300, 30), 2.0, 0.0, true, 0.0},
        {rl.NewVector2(900, 230), 2.0, 0.0, true, 0.0},
    }
}

func enemyCollides(position rl.Vector2, level Level) bool {
    enemyRect := Vector2ToRectangle(position, 20, 20)
    for _, room := range level.Rooms {
        if rl.CheckCollisionRecs(enemyRect, room.Bounds) {
            return true
        }
    }
    return false
}


// Weapon definitions

type Bullet struct {
	Position rl.Vector2
	Speed    float32
	Rotation float32
	Active   bool
    shotByEnemy bool // This flag indicates if the bullet was shot by an enemy
}

// Level definitions

type Rectangle struct {
	Bounds rl.Rectangle
	Color  rl.Color
}

type Level struct {
	Rooms []Rectangle
}

func Vector2ToRectangle(v rl.Vector2, width, height float32) rl.Rectangle {
	return rl.NewRectangle(v.X-width/2, v.Y-height/2, width, height)
}

func InitializeLevel() Level {
	return Level{
		Rooms: []Rectangle{
			// Outer boundary
			{rl.NewRectangle(0, 0, 4000, 30), rl.Gray},    // Top boundary
			{rl.NewRectangle(0, 0, 30, 3000), rl.Gray},    // Left boundary
			{rl.NewRectangle(3970, 0, 30, 3000), rl.Gray}, // Right boundary
			{rl.NewRectangle(0, 2970, 4000, 30), rl.Gray}, // Bottom boundary

			// Central room
			{rl.NewRectangle(1500, 1200, 1000, 600), rl.Gray},

			// Hallways
			{rl.NewRectangle(1200, 1400, 300, 100), rl.Gray},
			{rl.NewRectangle(2500, 1400, 300, 100), rl.Gray},
			{rl.NewRectangle(1800, 900, 400, 100), rl.Gray},
			{rl.NewRectangle(1800, 1900, 400, 100), rl.Gray},

			// Smaller rooms
			{rl.NewRectangle(800, 800, 400, 400), rl.Gray},
			{rl.NewRectangle(2800, 800, 400, 400), rl.Gray},
			{rl.NewRectangle(800, 1800, 400, 400), rl.Gray},
			{rl.NewRectangle(2800, 1800, 400, 400), rl.Gray},

			// Isolated rooms
			{rl.NewRectangle(500, 500, 200, 200), rl.Gray},
			{rl.NewRectangle(3300, 500, 200, 200), rl.Gray},
			{rl.NewRectangle(500, 2300, 200, 200), rl.Gray},
			{rl.NewRectangle(3300, 2300, 200, 200), rl.Gray},
		},
	}
}

func CheckCollision(player *Player, level Level) {
    futurePosition := player.Position

    if rl.IsKeyDown(rl.KeyW) {
        futurePosition.Y -= player.Speed
    }
    if !isColliding(futurePosition, level) {
        player.Position.Y = futurePosition.Y
    }

    futurePosition = player.Position
    if rl.IsKeyDown(rl.KeyS) {
        futurePosition.Y += player.Speed
    }
    if !isColliding(futurePosition, level) {
        player.Position.Y = futurePosition.Y
    }

    futurePosition = player.Position
    if rl.IsKeyDown(rl.KeyA) {
        futurePosition.X -= player.Speed
    }
    if !isColliding(futurePosition, level) {
        player.Position.X = futurePosition.X
    }

    futurePosition = player.Position
    if rl.IsKeyDown(rl.KeyD) {
        futurePosition.X += player.Speed
    }
    if !isColliding(futurePosition, level) {
        player.Position.X = futurePosition.X
    }
}

func isColliding(position rl.Vector2, level Level) bool {
    playerRect := Vector2ToRectangle(position, 20, 20)
    for _, room := range level.Rooms {
        if rl.CheckCollisionRecs(playerRect, room.Bounds) {
            return true
        }
    }
    return false
}

func DrawLevel(level Level) {
	for _, room := range level.Rooms {
		rl.DrawRectangleRec(room.Bounds, room.Color)
	}
}

func BulletCollidesWithWall(bullet Bullet, level Level) bool {
    bulletRect := Vector2ToRectangle(bullet.Position, 5, 5) // Assuming bullet radius is 5
    for _, room := range level.Rooms {
        if rl.CheckCollisionRecs(bulletRect, room.Bounds) {
            return true
        }
    }
    return false
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Top-Down Shooter")

    player := Player{
        Position: rl.NewVector2(2250, 1950), 
        Speed:    4.0,
    }

    enemies := InitializeEnemies()

    level := InitializeLevel()

	var bullets [10]Bullet

    camera := rl.NewCamera2D(rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)), player.Position, 0, 1)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Update
        mousePos := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera)
		
        // Update player
        player.Rotation = float32(math.Atan2(float64(mousePos.Y-player.Position.Y), float64(mousePos.X-player.Position.X))) * float32(180.0/math.Pi)

        // Update bullets (including player bullet collision with enemies)
        if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
            for i := range bullets {
                if !bullets[i].Active {
                    bullets[i].Position = player.Position
                    bullets[i].Rotation = player.Rotation
                    bullets[i].Speed = 15.0
                    bullets[i].Active = true
                    bullets[i].shotByEnemy = false // Set the flag to false when the bullet is shot by the player.
                    break
                }
            }
        }

        for i := range bullets {
            if bullets[i].Active {
                bullets[i].Position.X += bullets[i].Speed * float32(math.Cos(float64(bullets[i].Rotation)*math.Pi/180.0))
                bullets[i].Position.Y += bullets[i].Speed * float32(math.Sin(float64(bullets[i].Rotation)*math.Pi/180.0))
            }
        }

        // Update enemies
        for i := range enemies {
            if enemies[i].Alive {

                // Check if enemy collides with wall
                futurePosition := enemies[i].Position
                futurePosition.X += enemies[i].Speed * float32(math.Cos(float64(enemies[i].Rotation)*math.Pi/180.0))
                futurePosition.Y += enemies[i].Speed * float32(math.Sin(float64(enemies[i].Rotation)*math.Pi/180.0))
                if enemyCollides(futurePosition, level) {
                    enemies[i].Rotation += 180.0
                }

                if enemies[i].ShootCounter >= enemyShootDelay {
                    enemies[i].ShootCounter = 0
                    for j := range bullets {
                        if !bullets[j].Active {
                            bullets[j].Position.X = enemies[i].Position.X + 20 * float32(math.Cos(float64(enemies[i].Rotation)*math.Pi/180.0))
                            bullets[j].Position.Y = enemies[i].Position.Y + 20 * float32(math.Sin(float64(enemies[i].Rotation)*math.Pi/180.0))
                            bullets[j].Rotation = enemies[i].Rotation
                            bullets[j].Speed = 10.0
                            bullets[j].Active = true
                            bullets[j].shotByEnemy = true  // Set the flag to true when the bullet is shot by an enemy.
                            break
                        }
                    }
                }

                distance := math.Sqrt(math.Pow(float64(enemies[i].Position.X-player.Position.X), 2) + math.Pow(float64(enemies[i].Position.Y-player.Position.Y), 2))
                if distance < enemyRange {
                    enemies[i].Rotation = float32(math.Atan2(float64(player.Position.Y-enemies[i].Position.Y), float64(player.Position.X-enemies[i].Position.X))) * float32(180.0/math.Pi)
                    enemies[i].Position.X += enemies[i].Speed * float32(math.Cos(float64(enemies[i].Rotation)*math.Pi/180.0))
                    enemies[i].Position.Y += enemies[i].Speed * float32(math.Sin(float64(enemies[i].Rotation)*math.Pi/180.0))
                    
                    enemies[i].ShootCounter++
                    if enemies[i].ShootCounter >= enemyShootDelay {
                        enemies[i].ShootCounter = 0
                        for j := range bullets {
                            if !bullets[j].Active {
                                bullets[j].Position = enemies[i].Position
                                bullets[j].Rotation = enemies[i].Rotation
                                bullets[j].Speed = 10.0
                                bullets[j].Active = true
                                break
                            }
                        }
                    }
                }
                
                for j := range bullets {
                    if bullets[j].Active {
                        if rl.CheckCollisionCircles(bullets[j].Position, 5, enemies[i].Position, 20) {
                            enemies[i].Alive = false
                            bullets[j].Active = false
                            break
                        }
                    }
                }
            }
        }

        // Check for bullet-enemy collision
        for i := range enemies {
            if enemies[i].Alive {
                for j := range bullets {
                    if bullets[j].Active {
                        if !bullets[j].shotByEnemy && rl.CheckCollisionCircles(bullets[j].Position, 5, enemies[i].Position, 20) {
                            enemies[i].Alive = false
                            bullets[j].Active = false
                            break
                        }
                    }
                }
            }
        }

        // Check for bullet-wall collision
        for i := range bullets {
            if bullets[i].Active {
                if BulletCollidesWithWall(bullets[i], level) {
                    bullets[i].Active = false
                    continue
                }
            }

            // Deactivate bullets out of screen bounds
            if bullets[i].Position.X > 4000 || bullets[i].Position.X < 0 || bullets[i].Position.Y > 3000 || bullets[i].Position.Y < 0 {
                bullets[i].Active = false
            }
        }

        // Check collision (player with walls)
        CheckCollision(&player, level)

        // Update camera to follow the player
		camera.Target = player.Position

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

        rl.BeginMode2D(camera)

        DrawLevel(level)

		// Player
        rl.DrawCircleV(player.Position, 20, rl.Blue)
        rl.DrawCircleV(rl.NewVector2(player.Position.X+float32(20*math.Cos(float64(player.Rotation)*math.Pi/180.0)), player.Position.Y+float32(20*math.Sin(float64(player.Rotation)*math.Pi/180.0))), 5, rl.Black)

        // Draw enemies
        for i := range enemies {
            if enemies[i].Alive {
                rl.DrawCircleV(enemies[i].Position, 20, rl.Red)
            }
        }

        // Bullets
		for i := range bullets {
			if bullets[i].Active {
				rl.DrawCircleV(bullets[i].Position, 5, rl.Red)
			}
		}

        rl.EndMode2D() // End 2D mode

		rl.EndDrawing()
	}

	rl.CloseWindow()
}


