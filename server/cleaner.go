package server

import (
	"fmt"
	"lommix/wichtelbot/server/store"
	"time"
)

// Game and Session Garbage Collector
func (app *AppState) CleanupRoutine() {
	for {

		time.Sleep(time.Minute)

		// cleaning up any left over game sessions
		expiredSessions, err := store.FindExpiredParties(app.Db)
		if err != nil {
			panic(err)
		}

		if len(expiredSessions) > 0 {
			fmt.Printf("Cleaning %d expired sessions\n", len(expiredSessions))
			for _, session := range expiredSessions {
				err = store.DeleteUsersInParty(app.Db, session.Id)
				if err != nil {
					panic(err)
				}
				err = session.Delete(app.Db)
				if err != nil {
					panic(err)
				}
			}
		}

		// cleaning session memeory
		app.Sessions.CleanupExpired()
	}
}
