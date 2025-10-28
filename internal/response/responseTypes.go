package response

type StatusCode int

const StatusOk = StatusCode(200)

const StatusBadRequest = StatusCode(400)

const StatusInternalServerError = StatusCode(500)
