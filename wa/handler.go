package wa

import whatsapp "github.com/Rhymen/go-whatsapp"

type WAHandler struct{}

func (*WAHandler) HandleTextMessage(message whatsapp.TextMessage)         {}
func (*WAHandler) HandleImageMessage(message whatsapp.ImageMessage)       {}
func (*WAHandler) HandleVideoMessage(message whatsapp.VideoMessage)       {}
func (*WAHandler) HandleAudioMessage(message whatsapp.AudioMessage)       {}
func (*WAHandler) HandleDocumentMessage(message whatsapp.DocumentMessage) {}
func (*WAHandler) HandleJsonMessage(message string)                       {}
