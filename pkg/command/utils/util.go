package utils

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

func GenerateCSVFile(path string, rows [][]string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create csv file %s\n", err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(rows)
	csvWriter.Flush()
	return nil
}

func GenerateHTMLFile(sourceCSV string, targetHTML string) error {
	data, err := ioutil.ReadFile(sourceCSV)
	if err != nil {
		return fmt.Errorf("failed to read csv file %s", err)
	}
	const htmlTemplate = `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
		<title>Perf dashboard</title>
		<script src="https://cdn.jsdelivr.net/npm/jquery@3.5.1/dist/jquery.min.js"></script>
		<script src="https://cdn.jsdelivr.net/npm/chart.js@2.9.3/dist/Chart.min.js"></script>
		<script>
			var chartColors = {
				red: 'rgb(255, 99, 132)',
				orange: 'rgb(255, 159, 64)',
				yellow: 'rgb(255, 205, 86)',
				green: 'rgb(75, 192, 192)',
				blue: 'rgb(54, 162, 235)',
				purple: 'rgb(153, 102, 255)',
				grey: 'rgb(201, 203, 207)'
			}
			var colorNames = Object.keys(chartColors)
			function getChart(data, from, to) {
				var xAxesLable = ""
				var labels = []
				var datasets = []
				var table = ""
				var relArr = data.split("\n")
				if (!from || !to) {
					from = 1
					to = relArr.length - 1
					$("#range-from").val(from)
					$("#range-to").val(to)
				}
				if (!$.isEmptyObject(relArr) && relArr.length > 1) {
					for (var i = 0; i < relArr.length; i++) {
						var values = relArr[i];
						table += "<tr>"
						if (!$.isEmptyObject(values.trim())) {
							var objArr = values.trim().split(",");
							for (var j = 0; j < objArr.length; j++) {
								if (i == 0) {
									if (j == 0) {
										xAxesLable = objArr[j]
									} else {
										var colorName = colorNames[(j - 1) % colorNames.length];
										var newColor = chartColors[colorName];
										datasets[j - 1] = {
											label: objArr[j],
											backgroundColor: newColor,
											borderColor: newColor,
											fill: false,
											data: [],
											borderWidth: 1,
											pointRadius: 1
										}
									}
								} else {
									if (i >= from && i <= to) {
										if (j == 0) {
											labels[i - from] = objArr[j]
										} else {
											datasets[j - 1].data[i - from] = objArr[j]
										}
									}
								}
								if (j == 0) {
									if (i == 0) {
										table += "<th>index</th>"
									} else {
										if (i >= from && i <= to) {
											table += "<th>" + i + "</th>"
										} else {
											table += "<td>" + i + "</td>"
										}
									}
								}
								if (i == 0 || i >= from && i <= to) {
									table += "<th>" + objArr[j] + "</th>"
								} else {
									table += "<td>" + objArr[j] + "</td>"
								}
							}
						}
						table += "</tr>"
					}
				}
				var config = {
					type: 'line',
					data: {
						labels,
						datasets
					},
					options: {
						responsive: true,
						title: {
							display: false,
							text: 'TITLE'
						},
						tooltips: {
							mode: 'index',
							intersect: false,
						},
						hover: {
							mode: 'nearest',
							intersect: true
						},
						scales: {
							xAxes: [{
								display: true,
								scaleLabel: {
									display: true,
									labelString: xAxesLable
								}
							}],
							yAxes: [{
								display: true,
								scaleLabel: {
									display: true,
									labelString: ''
								}
							}]
						}
					}
				}
				return {
					config,
					table
				}
			}
			function jsReadFiles(files) {
				if (files.length) {
					var file = files[0];
					var reader = new FileReader();
					if (/text+/.test(file.type)) {
						reader.onload = function () {
							csvResult = this.result
							var chartInfo = getChart(csvResult)
							chart.data = chartInfo.config.data
							chart.update()
							$("#canvas-perf-detail-table").html("<tbody>" + chartInfo.table + "</tbody>")
						}
						reader.readAsText(file);
					} else {
						alert('Unsupported file')
					}
				}
			}
		</script>
		<style type="text/css">
			body {
				font-size: 14px;
			}
			table {
				margin-top: 20px;
				border-collapse: collapse;
			}
			table td,
			table th {
				border: 1px solid #dee2e6;
				text-align: center;
				padding: .75rem;
			}
			.perf-title {
				font-size: 16px;
				font-weight: bold;
				margin: 5px;
				margin-left: 0;
			}
			.perf-file {
				width: 45%;
				margin: 10px 20px 10px 20px;
				display: inline-block;
			}
			.perf-file input {
				width: 100%;
			}
			.perf-range {
				width: 45%;
				margin: 10px 20px 10px 20px;
				display: inline-block;
			}
			.perf-range input {
				width: unset;
				margin: 0 5px 0 5px;
			}
			.perf-container {
				width: 90%;
				margin: 20px;
				margin-top: 0px;
			}
		</style>
	</head>
	<body class="perf-detail-page">
		<div class="perf-file">
			<div class="perf-title">Data Source:</div>
			<input type="file" onchange="jsReadFiles(this.files)" />
		</div>
		<div class="perf-range">
			<div class="perf-title">Data Range:</div>
			From
			<input type="number" id="range-from" />
			To
			<input type="number" id="range-to" />
			<button type="button" class="btn btn-default" id="range-submit">Submit</button>
		</div>
		<div class="perf-container">
			<canvas id="canvas-perf-detail"></canvas>
			<table id="canvas-perf-detail-table"></table>
		</div>
		<script>
			var csvResult = "{{.Data}}"
			var chartInfo = getChart(csvResult)
			var ctx = $("#canvas-perf-detail")
			var chart = new Chart(ctx, chartInfo.config)
			$("#canvas-perf-detail-table").html("<tbody>" + chartInfo.table + "</tbody>")
			$("#range-submit").click(function () {
				var from = Number($("#range-from").val())
				var to = Number($("#range-to").val())
				if (isNaN(from) || isNaN(to) || from <= 0 || to <= 0) {
					alert("Range from & to must be numbers greater than 0")
					return
				}
				if (from > to) {
					alert("Range from must less or equal than to")
					return
				}
				var chartInfo = getChart(csvResult, from, to)
				chart.data = chartInfo.config.data
				chart.update()
				$("#canvas-perf-detail-table").html("<tbody>" + chartInfo.table + "</tbody>")
			})
		</script>
	</body>
	</html>
	`
	viewTemplate, err := template.New("chart").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse html template %s", err)
	}
	htmlFile, err := os.OpenFile(targetHTML, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open html file %s", err)
	}
	defer htmlFile.Close()
	return viewTemplate.Execute(htmlFile, map[string]interface{}{
		"Data": string(data),
	})
}
