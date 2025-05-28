package constants

const (
	TemplateTransactionEmail = `
		<html>
	<head><style>
	table { width: 100%; border-collapse: collapse; }
	th, td { padding: 8px; border: 1px solid #ccc; }
	th { background: #f5f5f5; text-align: left; }
	</style></head>
	<body>
		<p>Hai {{ .Name }},</p>
		<p>Pesanan Anda telah berhasil diproses. Berikut detailnya:</p>
		<table>
			<thead>
				<tr>
					<th>Product Name</th><th>Quantity</th><th>Product Price</th><th>Total Price</th>
				</tr>
			</thead>
			<tbody>
			{{- range .Items }}
				<tr>
					<td>{{ .ProductName }}</td>
					<td>{{ .Quantity }}</td>
					<td>Rp. {{ .ProductPrice }}</td>
					<td>Rp. {{ .TotalPrice }}</td>
				</tr>
			{{- end }}
			</tbody>
		</table>
		<p>
			<b>Total Quantity:</b> {{ .TotalQuantity }}<br/>
			<b>Total Discount:</b> Rp. {{ .Disc }}<br/>
			<b>Total Tax ({{ .TaxValue }}%):</b> Rp. {{ .TaxAmount }}<br/>
			<b>Total Product Amount:</b> Rp. {{ .TotalProductAmount }}<br/>
			<b>Total Transaction:</b> Rp. {{ .TotalTransactionAmount }}
		</p>
		<p>Terima kasih atas pembelian Anda!</p>
	</body>
	</html>
	`
)
