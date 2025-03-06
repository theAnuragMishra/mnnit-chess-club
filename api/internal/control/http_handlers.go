package control

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/notnil/chess"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

func (c *Controller) InitGame(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(event)
	var initGamePayload InitGamePayload
	err := json.NewDecoder(r.Body).Decode(&initGamePayload)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	session := r.Context().Value(auth.MiddlewareSentSession).(database.Session)

	log.Println("inside init game :)")

	var PendingUserID *int32
	var PendingUserName *string
	var baseTime time.Duration
	var increment time.Duration

	switch initGamePayload.TimerCode {
	case 1:
		PendingUserID = &c.GameManager.PendingUserID1
		PendingUserName = &c.GameManager.PendingUserName1
		baseTime = time.Minute
		increment = time.Duration(0)
	case 2:
		PendingUserID = &c.GameManager.PendingUserID2
		PendingUserName = &c.GameManager.PendingUserName2
		baseTime = time.Minute * 3
		increment = time.Second * 2
	case 3:
		PendingUserID = &c.GameManager.PendingUserID3
		PendingUserName = &c.GameManager.PendingUserName3
		baseTime = time.Minute * 10
		increment = time.Duration(0)
	default:
		PendingUserName = nil
		PendingUserID = nil

	}

	if PendingUserID == nil || PendingUserName == nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid timerCode")
		return
	}

	// if there's no pending game, create one, else, add the player to the game
	if *PendingUserID == 0 {
		log.Println("no pending game, creating...")
		*PendingUserName = initGamePayload.Username
		*PendingUserID = session.UserID

	} else {

		if *PendingUserName == initGamePayload.Username {
			utils.RespondWithError(w, http.StatusBadRequest, "You can't play both sides")
			return
		}
		createdGame := game.NewGame(baseTime, increment, *PendingUserID, session.UserID)

		id, err := c.Queries.CreateGame(context.Background(), database.CreateGameParams{
			BaseTime:      int32(baseTime.Seconds()),
			Increment:     int32(increment.Seconds()),
			WhiteID:       PendingUserID,
			BlackID:       &session.UserID,
			WhiteUsername: *PendingUserName,
			BlackUsername: initGamePayload.Username,
			Fen:           createdGame.Board.FEN(),
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		createdGame.ID = id
		c.GameManager.Games = append(c.GameManager.Games, createdGame)

		thisClient := c.SocketManager.FindClientByUserID(session.UserID)
		otherClient := c.SocketManager.FindClientByUserID(*PendingUserID)
		if thisClient == nil || otherClient == nil {
			log.Println("player not found")
			return
		}

		*PendingUserID = 0
		*PendingUserName = ""
		payload := map[string]interface{}{
			"GameID": createdGame.ID,
		}
		rawPayload, err := json.Marshal(payload)
		if err != nil {
			log.Println("error marshalling new createdGame payload")
			return
		}
		e := socket.Event{
			Type:    "Init_Game",
			Payload: json.RawMessage(rawPayload),
		}

		thisClient.Send(e)
		otherClient.Send(e)

	}

	// fmt.Println(c.GameManager)
}

func (c *Controller) WriteProfileInfo(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	profileInfo, err := c.Queries.GetPlayerGames(r.Context(), username)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	}

	utils.RespondWithJSON(w, http.StatusOK, profileInfo)
}

func (c *Controller) WriteGameInfo(w http.ResponseWriter, r *http.Request) {
	gameIDStr, err := strconv.ParseInt(chi.URLParam(r, "gameID"), 10, 32)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID")
	}
	gameID := int32(gameIDStr)

	// fmt.Println(gameID)
	foundGame, err := c.Queries.GetGameInfo(r.Context(), gameID)
	if err != nil {
		log.Println(err)

		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID")
	}
	moves, err := c.Queries.GetGameMoves(r.Context(), gameID)
	if err != nil {
		log.Println(err)
	}

	var serverGame *game.Game
	var ongoing bool
	var timeBlack int32
	var timeWhite int32
	for _, g := range c.GameManager.Games {
		if g.ID == gameID {
			serverGame = g
		}
	}
	if serverGame == nil || foundGame.Result != "ongoing" {
		ongoing = false
		if foundGame.EndTimeLeftWhite != nil {
			timeWhite = *foundGame.EndTimeLeftWhite
		}
		if foundGame.EndTimeLeftBlack != nil {
			timeBlack = *foundGame.EndTimeLeftBlack
		}

	} else {
		timePassed := time.Since(serverGame.LastMoveTime)
		if serverGame.Board.Position().Turn() == chess.White {
			serverGame.TimeWhite = max(serverGame.TimeWhite-timePassed, 0)
		} else {
			serverGame.TimeBlack = max(serverGame.TimeBlack-timePassed, 0)
		}
		serverGame.LastMoveTime = time.Now()
		ongoing = true
		timeBlack = int32(serverGame.TimeBlack.Seconds())
		timeWhite = int32(serverGame.TimeWhite.Seconds())
	}

	// fmt.Println(serverGame.Result)

	response := GameResponse{
		Game:      foundGame,
		Moves:     moves,
		Ongoing:   ongoing,
		TimeBlack: timeBlack,
		TimeWhite: timeWhite,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}
