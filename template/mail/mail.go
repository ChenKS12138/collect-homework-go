package mail

// InvitationCodeTemplate invitation code template
const InvitationCodeTemplate = 
`<div>
	<p style="text-align:center;font-size:20px">这是来自 {{.Email}} 的申请，请审核</p>
	<p style="text-align:center;font-size:30px">邀请码： <a style="color:red">{{.Code}}</a> </p>
	<p style="text-align:center;font-size:20px;text-decoration:underline">若信任该申请，请将邀请码分享给TA</p>
	<p style="text-align:center;font-size:13px">时间: <span style="text-decoration: underline">{{.Time}}</span></p>
</div>`

// SubmissionTemplate submissionTemplate
const SubmissionTemplate = 
`<div>
	<p style="text-align:center;font-size:20px">作业项目: <span style="text-decoration: underline">{{.ProjectName}}</span></p>
	<p style="text-align:center;font-size:30px">{{.Status}} <span style="color:red">{{.FileName}}</span></p>
	<p style="text-align:center;font-size:13px">时间: <span style="text-decoration: underline">{{.Time}}</span></p>
	<p style="text-align:center;color:gray;font-size:13px">IP地址: {{.IP}}</p>
</div>`