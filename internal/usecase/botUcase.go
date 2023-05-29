package usecase

import (
	"botserver/internal/asr"
	"botserver/internal/model"
	"botserver/internal/nlp"
	"botserver/pkg/natsclient"

	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
)

type BotUsecase struct {
	natsClient *natsclient.Client
	asrPool    *asr.ASRPool
	nlp        *nlp.NLP
}

func NewBotUsecase(natsClient *natsclient.Client) *BotUsecase {
	asrPool := asr.NewASRPool(3)
	nlpComponent := nlp.NewNLP()

	return &BotUsecase{
		natsClient: natsClient,
		asrPool:    asrPool,
		nlp:        nlpComponent,
	}
}

func (u *BotUsecase) StartBotServer() error {
	// обрабатывать входящие сообщения
	subInput, err := u.natsClient.Subscribe("player.input", u.handlePlayerInput)
	if err != nil {
		return err
	}

	subOutput, err := u.natsClient.Subscribe("player.output", u.handlePlayerOutput)
	if err != nil {
		// Отменить подписку на тему player.input, если произошла ошибка
		subInput.Unsubscribe()
		return err
	}

	done := make(chan bool)

	go func() {
		subInput.Unsubscribe()
		subOutput.Unsubscribe()
		done <- true
	}()

	<-done

	u.asrPool.Stop()

	return nil
}

func (u *BotUsecase) handlePlayerInput(msg *nats.Msg) {
	audioChunk := msg.Data

	// Обработка аудиофрагмента
	transcription, err := u.asrPool.ProcessAudioChunk(audioChunk)
	if err != nil {
		log.Printf("Failed to process audio chunk: %v", err)
		return
	}

	// Обработка транскрипции
	result, err := u.nlp.ProcessText(transcription)
	if err != nil {
		log.Printf("Failed to process transcription: %v", err)
		return
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Printf("Failed to marshal result to JSON: %v", err)
		return
	}

	if err := u.natsClient.Publish("ui.output", resultJSON); err != nil {
		log.Printf("Failed to publish result: %v", err)
		return
	}
}

func (u *BotUsecase) handlePlayerOutput(msg *nats.Msg) {
	// Обработать сообщение полученное от NATS
	var playerStatus model.PlayerStatus
	err := json.Unmarshal(msg.Data, &playerStatus)
	if err != nil {
		// Handle the JSON parsing error
		log.Println("Failed to parse player output message:", err)
		return
	}

	// Handle the player status update
	//u.updateUIWithPlayerStatus(playerStatus)

	// Send the player status update to the frontend
	//u.sendPlayerStatusToFrontend(playerStatus)
}