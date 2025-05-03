package db

type DocumentDAO struct {
	cm *ConnectionManager
}

func NewDocumentDAO(cm *ConnectionManager) *DocumentDAO {
	return &DocumentDAO{
		cm: cm,
	}
}
