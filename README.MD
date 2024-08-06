Run Server Program first.

go run server/server.go

Run the Client Program with the command line arguments.

Example:

go run client/client.go purchase London France John Doe john.doe@example.com

go run client/client.go get_receipt rec-1

go run client/client.go view_users SectionA

go run client/client.go modify_seat john.doe@example.com NewSeat-1

go run client/client.go remove_user john.doe@example.com
