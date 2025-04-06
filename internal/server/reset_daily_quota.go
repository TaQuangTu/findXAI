package server

import (
	"findx/config"
	"findx/internal/lockdb"
	"log"
	"time"
)

// StartDailyResetTask launches the background process to reset daily counts at midnight UTC-7
func StartDailyResetTask(cfg *config.Config, lockDb lockdb.ILockDb, searchServer *SearchServer) {
	utcMinus7 := getUTCMinus7Location()

	for {
		performDailyReset(cfg, lockDb, searchServer)
		sleepUntilNextReset(utcMinus7)
	}
}

// getUTCMinus7Location returns the timezone location for UTC-7
func getUTCMinus7Location() *time.Location {
	utcMinus7, err := time.LoadLocation("America/Los_Angeles") // Pacific Time approximates UTC-7/UTC-8
	if err != nil {
		log.Printf("Failed to load timezone: %v, falling back to UTC-7 offset", err)
		utcMinus7 = time.FixedZone("UTC-7", -7*60*60)
	}
	return utcMinus7
}

// sleepUntilNextReset calculates and waits until the next reset time
func sleepUntilNextReset(utcMinus7 *time.Location) {
	// Calculate time until next midnight in UTC-7
	now := time.Now().In(utcMinus7)
	thisMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, utcMinus7)
	nextMidnight := thisMidnight.Add(24 * time.Hour)

	// Calculate duration until next reset
	sleepDuration := nextMidnight.Sub(now)
	log.Printf("Scheduled next daily count reset in %v (at %v UTC-7)", sleepDuration, thisMidnight)

	// Sleep until the next reset time
	time.Sleep(sleepDuration)
}

// performDailyReset attempts to acquire a lock and reset the daily counts
func performDailyReset(cfg *config.Config, lockDb lockdb.ILockDb, searchServer *SearchServer) {

	// TODO: Try to acquire lock to ensure only one instance resets the counts
	// goCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// ourLock, err := lockDb.LockSimple(goCtx, "search:daily_count_reset")
	// if err != nil {
	// 	log.Printf("Failed to acquire lock for daily count reset: %v", err)
	// 	return
	// }
	// defer ourLock.Unlock()

	// Reset daily counts
	searchServer.KeyManager.ResetDailyCounts(cfg.MAX_REQUEST_PER_DAY)
	log.Printf("Daily API counts reset at %v (midnight UTC-7)", time.Now().UTC())
}
