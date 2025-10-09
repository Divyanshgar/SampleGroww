# TODO: Update main.go to initialize PDF service and pass to notification handler

- [x] Update main.go: Initialize PDF service and pass to routes.SetupRoutes
- [x] Update routes/routes.go: Modify SetupRoutes to accept pdfService and pass to NewNotificationHandler
- [x] Update handlers/notification_handler.go: Modify NewNotificationHandler to accept pdfService parameter
