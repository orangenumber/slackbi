package slackbi

import (
	"fmt"
)

// msg/msge is from github.com/gonyyi/tmp/anerr
type msge string

func (e msge) String() string {
	return string(e)
}

type msg string

func (m msg) String() string {
	return string(m)
}
func (m msg) Error() string { // to meet error interface
	return string(m)
}
func (m msg) Format(a ...interface{}) string {
	return fmt.Sprintf(m.String(), a...)
}
func (m msg) Errorf(a ...interface{}) msge {
	return msge(m.Format(a...))
}

const (
	// Message
	M_SYSF_UPDATED msg = "sysf has been updated"
	M_SYSF_INVALID msg = "invalid sysf function"
	M_HTTP_RESP_OK msg = "ok"
	M_MOD_ERR_001  msg = "ERR-001: problem with the module."

	// Message Formatted
	MF_HTTP_SERVING_SAddr_SPort_SPath msg = "Serving HTTP <%s%s%s>"
	MF_SBI_SVersion                   msg = "SlackBotInterface %s"
	MF_SBI_CREATING_SName_SVersion    msg = "Creating %s: %s"
	MF_SLACK_BOT_CHALLENGE_SChallenge msg = "Received a slack bot challenge request; challenge=%s"

	MF_MSG_UNMARSHAL_FAILED_SErr                   msg = "MSG received but can't unmarshal, err=%s"
	MF_MSG_UNMARSHAL_FAILED_SErr_SData             msg = "MSG received but can't unmarshal, err=%s, data=%s"
	MF_MSG_RECEIVED_SFrom_SMsg_SThread             msg = "MSG msg_received, From=%s, msg=%s, Thread=%s"
	MF_MSG_RECEIVED_IGNORE_SThread_SBotID_SSubType msg = "MSG msg_received but ignored, TS=%s, BotID=%s, SubType=%s"
	MF_MSG_READ_FAILED_SErr                        msg = "MSG received but failed to read, err=%s"
	MF_MSG_RECEIVED_SData                          msg = "Received a call, data=%s"

	MF_MSG_RESP_EMPTY_SThread            msg = "response doesn't have any message, incomingMsg.TS=%s"
	MF_MSG_RESP_FILE_DELETE_FAILED_SFile msg = "uploaded but failed to delete, file=%s"
	MF_MSG_RESP_FILE_UPLOADED_SFile      msg = "uploaded and deleted, file=%s"

	MF_MSG_SENT_BUT_ERROR_SErr msg = "sent request, error=%s"
	MF_MSG_SENT_OK_SData       msg = "sent request, received data=%s"

	MF_FILE_SENT_SData msg = "sent request, received data=%s"

	MF_CONF_OVERRIDE_SName_SVal_SNewVal msg = "Validation failed, override config, for <%s> <%s> -> <%s>"
	MF_HTTP_RESP_UNEXPECTED_SErr        msg = "Unexpected: %s"

	MF_MOD_DELAYED_RESP_SAvgSec                         msg = "Wait.. average runtime for this module is %d seconds."
	MF_MOD_OUTPUT_SData                                 msg = "output=%s"
	MF_MOD_ERR_001_SErr                                 msg = "ERR-001: err=%s"
	MF_MOD_ERR_002_SErr                                 msg = "ERR-002: *MsgIncoming.ResponseText(), err=%s"
	MF_MOD_ERR_JSON_SErr                                msg = "in.JSON() failed, err= %s"
	MF_MOD_EXEC_V2_FAILED_SErr_SOutput                  msg = "module.ExecV2() failed, err=%s, output=%s"
	MF_MOD_EXEC_V2_OK_SData                             msg = "received from ExecV2(), data=%s"
	MF_MOD_EXEC_V2_UNMARHSAL_FAILED_SErr_SOutput        msg = "Unmarshal output to MsgOutgoing failed, err=%s, output=%s"
	MF_MOD_EXEC_V2_RESP_FAILED_SErr_SOutput             msg = "in.Response() failed, err=%s, output=%s"
	MF_MOD_EXEC_V2_STDOUT_EMPTY_IStdoutSize_IStderrSize msg = "no stdout received, stdout.size=%d, stderr.size=%d"

	MF_CMD_UNKNOWN_SCommand msg = "Sorry. <%s> is an unknown command or a module"
)
