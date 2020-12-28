package slackbi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// =====================================================================================================================
// HANDLING COMMANDS
// =====================================================================================================================

func (b *SBI) command(in MsgIncoming) {
	txt := in.Text()

	oops := func(err error) {
		b.logger.Errorf("ERR-001: err=%s", err.Error())
		if err = in.ResponseText(b, false, "ERR-001: problem with the module."); err != nil {
			b.logger.Errorf("ERR-002: *MsgIncoming.ResponseText(), err=%s", err.Error())
		}
	}

	// Split first word and check if exists in modules
	moduleName := strings.Split(txt, " ")[0]
	if b.modules.IsExist(moduleName) {
		module, err := b.modules.Get(moduleName)
		if err != nil {
			oops(err)
			return
		}

		outmsg := in.OutgoingMsg(b)

		if module.Module.AvgRuntimeSec == 0 {
			module.Module.AvgRuntimeSec = 1
		}
		if module.Module.AvgRuntimeSec > 5 {
			in.ResponseText(b, false, "Wait.. average runtime for this module is "+strconv.Itoa(module.Module.AvgRuntimeSec)+" seconds.")
		}

		switch module.Module.InterfaceVersion {
		case 1:
			output, outputError, err := module.ExecV1(in.Event.User, txt)
			if err != nil {
				oops(err)
				return
			}
			if len(output) > 0 {
				outmsg.Blocks.AddMarkdown(string(output))
			}
			if len(outputError) > 0 {
				outmsg.Blocks.AddDivider()
				tmpME := MsgElement{}
				tmpME.AsMarkdown("```" + string(outputError) + "```")
				outmsg.Blocks.AddContext(tmpME)
				outmsg.Send(b)
				b.logger.Tracef("output=%s", string(outmsg.JSON()))
			}

		default:
			// module 2 is default
			jIn, err := in.JSON()
			if err != nil {
				b.logger.Tracef("in.JSON() failed, err=", err.Error())
				oops(err)
				return
			}
			output, outputError, err := module.ExecV2(jIn)
			if err != nil {
				b.logger.Tracef("module.ExecV2() failed, err=%s, output=%s", err.Error(), string(output))
				oops(err)
				return
			}

			b.logger.Tracef("received from ExecV2(), data=%s", string(output))

			if len(output) > 0 {
				var tmpMO MsgOutgoing
				if err := json.Unmarshal(output, &tmpMO); err != nil {
					b.logger.Tracef("Unmarshal output to MsgOutgoing failed, err=%s, output=%s", err.Error(), string(output))
					oops(err)
					return
				}
				if err := in.Response(b, moduleName, &tmpMO); err != nil {
					b.logger.Tracef("in.Response() failed, err=%s, output=%s", err.Error(), string(output))
					oops(err)
					return
				}
			} else {
				b.logger.Tracef("no stdout received, stdout.size=%d, stderr.size=%d", len(output), len(outputError))
			}

			if len(outputError) > 0 {
				in.ResponseMarkdown(b, false, "*Error:* ```"+string(outputError)+"```")
			}
		}
	} else {
		b.Help(&in)
	}
}


func (b *SBI) Help(in *MsgIncoming) {
	in.ResponseText(b, false, fmt.Sprintf("Sorry. <%s> is an unknown command or a module", in.Text()))
}
