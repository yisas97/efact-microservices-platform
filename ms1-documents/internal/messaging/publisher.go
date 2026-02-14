package messaging

type Publisher interface {
	PublishDocumentCreated(documentID, uuid string) error
}
