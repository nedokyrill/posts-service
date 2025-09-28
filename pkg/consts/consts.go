package consts

import "time"

const PageSize = 20
const ContentMaxLen = 2000

const InitPostsSizeInMem = 200
const InitCommentsSizeInMem = 50

const BadRequestType = "Bad Request"
const InternalServerErrorType = "Internal Server Error"

const PgxTimeout = 5 * time.Second
