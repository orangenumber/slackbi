package slackbi

import (
	"encoding/json"
	"fmt"
	"strings"
)

// =====================================================================================================================
// HANDLING COMMANDS
// =====================================================================================================================

func (b *SBI) command(in MsgIncoming) {
	txt := in.Text()

	oops := func(err error) {
		b.logger.Errorf(MF_MOD_ERR_001_SErr.Format(err.Error()))
		if err = in.ResponseText(b, false, m_mod_err_001.String()); err != nil {
			b.logger.Errorf(MF_MOD_ERR_002_SErr.Format(err.Error()))
		}
	}

	// Split first word and check if exists in modules
	moduleName := strings.Split(txt, " ")[0]
	if module, err := b.modules.Get(moduleName); err == nil {
		outmsg := in.OutgoingMsg(b)
		outmsg.Custom.ReplyInThread = false

		if module.Module.AvgRuntimeSec == 0 {
			module.Module.AvgRuntimeSec = 1
		}
		if module.Module.AvgRuntimeSec > 5 {
			in.ResponseText(b, false,
				MF_MOD_DELAYED_RESP_SAvgSec.Format(module.Module.AvgRuntimeSec))
		}

		switch module.Module.InterfaceVersion {
		case 1:
			output, outputError, err := module.ExecV1(in.Event.User, txt)
			if err != nil {
				oops(err)
				return
			}
			if len(output) > 0 {
				outmsg.Blocks.AddMarkdown("```" + string(output) + "```")
			}
			if len(outputError) > 0 {
				outmsg.Blocks.AddDivider()
				tmpME := MsgElement{}
				tmpME.AsMarkdown("```" + string(outputError) + "```")
				outmsg.Blocks.AddContext(tmpME)
			}
			outmsg.Send(b)
			b.logger.Tracef(MF_MOD_OUTPUT_SData.Format(string(outmsg.JSON())))

		default:
			// module 2 is default
			jIn, err := in.JSON()
			if err != nil {
				b.logger.Tracef(MF_MOD_ERR_JSON_SErr.Format(err.Error()))
				oops(err)
				return
			}
			output, outputError, err := module.ExecV2(jIn)
			if err != nil {
				b.logger.Tracef(MF_MOD_EXEC_V2_FAILED_SErr_SOutput.Format(err.Error(), string(output)))
				oops(err)
				return
			}

			b.logger.Tracef(MF_MOD_EXEC_V2_OK_SData.Format(string(output)))

			if len(output) > 0 {
				var tmpMO MsgOutgoing
				if err := json.Unmarshal(output, &tmpMO); err != nil {
					b.logger.Tracef(MF_MOD_EXEC_V2_UNMARHSAL_FAILED_SErr_SOutput.Format(err.Error(), string(output)))
					oops(err)
					return
				}
				if err := in.Response(b, moduleName, &tmpMO); err != nil {
					b.logger.Tracef(MF_MOD_EXEC_V2_RESP_FAILED_SErr_SOutput.Format(err.Error(), string(output)))
					oops(err)
					return
				}
			} else {
				b.logger.Tracef(MF_MOD_EXEC_V2_STDOUT_EMPTY_IStdoutSize_IStderrSize.Format(len(output), len(outputError)))
			}

			if len(outputError) > 0 {
				in.ResponseMarkdown(b, false, "*Error:* ```"+string(outputError)+"```")
			}
		}
	} else {
		in.ResponseText(b, false, fmt.Sprintf("Module <%s> not found", moduleName))
	}
}

func (b *SBI) Help(in *MsgIncoming) {
	in.ResponseText(b, false, MF_CMD_UNKNOWN_SCommand.Format(in.Text()))
}
