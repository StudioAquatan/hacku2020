package character

import (
	"log"
	"math/rand"
	"time"
)

type Info struct {
	Name    string
	Icon    string
	Message []string
}

type MessageInfo struct {
	Name    string
	Icon    string
	Message string
}

func (ci *Info) removeMessage(index int) {
	var ms []string
	for i, v := range ci.Message {
		if i == index {
			continue
		}
		ms = append(ms, v)
	}
	ci.Message = ms
}

func removeCharacter(cis *[]Info, index int) []Info {
	var newCis []Info
	for i, v := range *cis {
		if i == index {
			continue
		}
		newCis = append(newCis, v)
	}
	return newCis
}

func CreateMessageInfoByRandom(cis []Info, messageNum int) *[]MessageInfo {
	var mis []MessageInfo

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < messageNum; i++ {
		characterIndex := rand.Intn(len(cis))
		messageIndex := rand.Intn(len(cis[characterIndex].Message))
		mi := MessageInfo{
			Name:    cis[characterIndex].Name,
			Icon:    cis[characterIndex].Icon,
			Message: cis[characterIndex].Message[messageIndex],
		}
		mis = append(mis, mi)
		cis[characterIndex].removeMessage(messageIndex)
		if len(cis[characterIndex].Message) == 0 {
			cis = removeCharacter(&cis, characterIndex)
		}
		if len(cis) == 0 {
			log.Println("[INFO] No left message")
			break
		}
	}

	return &mis
}
