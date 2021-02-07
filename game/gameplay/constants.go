package gameplay

import "time"

const championRadius = 75
const championStartX = 1500
const championStartY = 1500
const championSpeed = 700
const championMaxHealth = 100

const dashCooldown = time.Second * 2
const dashDuration = time.Second / 5
const dashSpeedMultiplier = 3
const shootCooldown = time.Second / 10

const projectileDamage = 10
const projectileSpeed = 1600
const projectileRange = 2000
const projectileRadius = 25
