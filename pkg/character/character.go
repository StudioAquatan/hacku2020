package character

import (
	"log"
	"math/rand"
	"time"
)

type Info struct {
	Name                   string
	Icon                   string
	EncourageMessages      []string
	CongratulatoryMessages []string
}

type MessageInfo struct {
	Name    string
	Icon    string
	Message string
}

func (ci *Info) removeMessage(index int, oinori bool) {
	var newMs []string
	var ms []string

	if oinori {
		ms = ci.EncourageMessages
	} else {
		ms = ci.CongratulatoryMessages
	}

	for i, v := range ms {
		if i == index {
			continue
		}
		newMs = append(newMs, v)
	}

	if oinori {
		ci.EncourageMessages = newMs
	} else {
		ci.CongratulatoryMessages = newMs
	}
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

func CreateMessageInfoByRandom(cis *[]Info, messageNum int, oinori bool) *[]MessageInfo {
	var mis []MessageInfo

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < messageNum; i++ {
		if oinori {
			characterIndex := rand.Intn(len(cis))
			messageIndex := rand.Intn(len(cis[characterIndex].EncourageMessages))
			mi := MessageInfo{
				Name:    cis[characterIndex].Name,
				Icon:    cis[characterIndex].Icon,
				Message: cis[characterIndex].EncourageMessages[messageIndex],
			}
			mis = append(mis, mi)
			cis[characterIndex].removeMessage(messageIndex, oinori)
			if len(cis[characterIndex].EncourageMessages) == 0 {
				cis = removeCharacter(&cis, characterIndex)
			}
			if len(cis) == 0 {
				log.Println("[INFO] No left message")
				break
			}
		} else {
			characterIndex := rand.Intn(len(cis))
			messageIndex := rand.Intn(len(cis[characterIndex].CongratulatoryMessages))
			mi := MessageInfo{
				Name:    cis[characterIndex].Name,
				Icon:    cis[characterIndex].Icon,
				Message: cis[characterIndex].CongratulatoryMessages[messageIndex],
			}
			mis = append(mis, mi)
			cis[characterIndex].removeMessage(messageIndex, oinori)
			if len(cis[characterIndex].CongratulatoryMessages) == 0 {
				cis = removeCharacter(&cis, characterIndex)
			}
			if len(cis) == 0 {
				log.Println("[INFO] No left message")
				break
			}
		}
	}

	return &mis
}
