From: Anygonow <{{ .EmailSender }}>
Subject: {{ .Subject }}
To: {{ .UserEmail }}
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="{{ .Boundary }}"

--{{ .Boundary }}
Content-Type: text/plain

Hi {{ .UserEmail }},

{{ .Text1 }}
{{ .Text2 }}
{{ .Link }}

For further questions, please contact: {{ .CompanyEmail }}

--{{ .Boundary }}
Content-Type: text/html

<!DOCTYPE html>
<html>

<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>AnyGoNow</title>
</head>

<body style="background-color: #ffffff">
	<div style="background-color: #f6f6f6; max-width: 600px; margin: 0 auto;">
		<table style="padding-left: 20px; padding-top: 30px">
			<tr>
				<td>
					<img src="https://d26p4pe0nfrg62.cloudfront.net/email/icon70x70.png"> </img>
				</td>
			</tr>
			<tr>
				<td>
					<div>Hi <span href="mailto:{{ .UserEmail }}">{{ .UserEmail }}</span>,</div>
				</td>
			</tr>
		</table>
		<table style="padding-left: 20px; padding-top: 10px">
			<tr>
				<td>
					<div>
            {{ .Text1 }}
					</div>
				</td>
			</tr>
			<tr>
				<td>
					<div>
            {{ .Text2 }}
					</div>
				</td>
			</tr>
			<tr>
				<td style="padding-top: 10px;">
						<div style="background-color: #ff511a;
						width: 100px;
						height: 30px;
						border-radius: 4px;
						text-align: center;" >
							<a style="text-align: center;
							display: inline-block;
							vertical-align: sub;
							text-decoration: none !important;
							color: black;
							margin: 0 auto" href="{{.Link }}">{{ .ButtonText }}</a>
						</div>
					
				</td>
			</tr>
		</table>
		<table style="padding-left: 20px; padding-top: 20px; padding-bottom: 30px;">
			<tr>
				<td>
					<div>
						For further questions, please contact: <span><a href="mailto:{{ .CompanyEmail }}">{{ .CompanyEmail }}</a></span>
					</div>
				</td>
			</tr>
		</table>
	</div>

</body>

</html>
--{{ .Boundary }}--