package slackbi

import (
	"fmt"
)

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
func (m msg) Errorf(a ...interface{}) msg {
	return msg(m.Format(a...))
}

const (
	// Message
	m_http_resp_ok msg = "ok"
	m_mod_err_001  msg = "ERR-001: problem with the module."
	m_config_empty msg = "config empty"

	mf_config_read_err_sErr msg = "cannot open config, err=%s"
	m_sbi_ptr_error         msg = "sbi not initiated correctly"

	// Message Formatted
	mf_app_sName_sVersion             msg = "%s: %s"
	mf_sbi_sVersion                   msg = "SlackBotInterface: %s"
	mf_http_serving_sAddr_sPort_sPath msg = "Serving HTTP <%s%s%s>"
	mf_slack_bot_challenge_sChallenge msg = "INCOMING, a slack bot challenge; challenge=%s"

	mf_msg_unmarshal_failed_sErr                   msg = "INCOMING, can't unmarshal, err=%s"
	mf_msg_unmarshal_failed_sErr_sData             msg = "INCOMING, can't unmarshal, err=%s, data=%s"
	mf_msg_received_sFrom_sMsg_sThread             msg = "INCOMING, From=%s, msg=%s, Thread=%s"
	mf_msg_received_ignore_sThread_sBotID_sSubtype msg = "INCOMING, ignored, TS=%s, BotID=%s, SubType=%s"
	mf_msg_read_failed_sErr                        msg = "INCOMING, failed to read, err=%s"
	mf_msg_received_sData                          msg = "INCOMING, data=%s"

	MF_MSG_RESP_EMPTY_SThread            msg = "OUTGOING, no response back, receivedTS=%s"
	MF_MSG_RESP_FILE_DELETE_FAILED_SFile msg = "OUTGOING, uploaded but failed to delete, file=%s"
	MF_MSG_RESP_FILE_UPLOADED_SFile      msg = "OUTGOING, deleted, file=%s"

	MF_MSG_SENT_BUT_ERROR_SErr msg = "OUTGOING, error=%s"
	MF_MSG_SENT_OK_SData       msg = "OUTGOING, OK, received resp=%s"

	MF_FILE_SENT_SData msg = "OUTGOING, received resp=%s"

	MF_CONF_OVERRIDE_SName_SVal_SNewVal msg = "Validation failed, override config, for <%s> <%s> -> <%s>"
	MF_HTTP_RESP_UNEXPECTED_SErr        msg = "Unexpected: %s"

	MF_MOD_DELAYED_RESP_SAvgSec                         msg = "Wait.. average runtime for this module is %d seconds."
	MF_MOD_OUTPUT_SData                                 msg = "output=%s"
	mf_mod_not_found_sName                              msg = "module %s not found"
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
