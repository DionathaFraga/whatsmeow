package whatsmeow

import (
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

// ButtonQuickReply representa um botão de resposta rápida
type ButtonQuickReply struct {
	DisplayText string
	ID          string
}

// NativeFlowButton representa um botão de flow
type NativeFlowButton struct {
	Name       string
	ButtonText string
}

// SendQuickReplyButtons envia mensagem com botões de resposta rápida
func (cli *Client) SendQuickReplyButtons(jid string, body string, footer string, buttons []ButtonQuickReply) error {
	recipient, err := cli.ParseJID(jid)
	if err != nil {
		return err
	}

	// Cria os botões
	var protoButtons []*waProto.Message_ButtonsMessage_Button
	for _, btn := range buttons {
		protoButtons = append(protoButtons, &waProto.Message_ButtonsMessage_Button{
			ButtonId:       &btn.ID,
			ButtonText:     &waProto.Message_ButtonsMessage_Button_ButtonText{DisplayText: &btn.DisplayText},
			Type:           waProto.Message_ButtonsMessage_Button_QUICK_REPLY.Enum(),
		})
	}

	// Cria a mensagem
	msg := &waProto.Message{
		ButtonsMessage: &waProto.Message_ButtonsMessage{
			ContentText: &body,
			FooterText:  &footer,
			Buttons:     protoButtons,
			HeaderType:  waProto.Message_ButtonsMessage_EMPTY.Enum(),
		},
	}

	_, err = cli.SendMessage(cli.Store.Context, recipient, msg)
	return err
}

// SendListMessage envia mensagem com lista de opções
func (cli *Client) SendListMessage(jid string, body string, footer string, title string, buttonText string, sections []ListSection) error {
	recipient, err := cli.ParseJID(jid)
	if err != nil {
		return err
	}

	// Converte sections para proto
	var protoSections []*waProto.Message_ListMessage_Section
	for _, section := range sections {
		var rows []*waProto.Message_ListMessage_Row
		for _, row := range section.Rows {
			rows = append(rows, &waProto.Message_ListMessage_Row{
				RowId:       &row.ID,
				Title:       &row.Title,
				Description: &row.Description,
			})
		}
		protoSections = append(protoSections, &waProto.Message_ListMessage_Section{
			Title: &section.Title,
			Rows:  rows,
		})
	}

	msg := &waProto.Message{
		ListMessage: &waProto.Message_ListMessage{
			Title:       &title,
			Description: &body,
			FooterText:  &footer,
			ButtonText:  &buttonText,
			ListType:    waProto.Message_ListMessage_SINGLE_SELECT.Enum(),
			Sections:    protoSections,
		},
	}

	_, err = cli.SendMessage(cli.Store.Context, recipient, msg)
	return err
}

// ListSection representa uma seção de lista
type ListSection struct {
	Title string
	Rows  []ListRow
}

// ListRow representa uma linha de lista
type ListRow struct {
	ID          string
	Title       string
	Description string
}
