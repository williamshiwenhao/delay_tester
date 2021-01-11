package main

func main() {
	receiverConn := CreateClient(Config.ReceiverAddress)
	senderConn := CreateClient(Config.SenderAddress)
	receiver := NewReceiver(receiverConn)
	sender := NewSender(senderConn)
	go receiver.Run()
	sender.Run()
}
