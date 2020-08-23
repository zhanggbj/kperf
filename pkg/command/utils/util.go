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
		<script src="https://cdn.jsdelivr.net/npm/echarts/dist/echarts-en.min.js"></script>
		<script>
			const legendIcons = [
				'path://M512 512m-448 0a448 448 0 1 0 896 0 448 448 0 1 0-896 0Z',
				'path://M442.7 145.8L66.1 798.2c-30.8 53.3 7.7 120 69.3 120h753.3c61.6 0 100.1-66.7 69.3-120L581.3 145.8c-30.8-53.3-107.8-53.3-138.6 0z',
				'path://M760 960H264c-110.5 0-200-89.5-200-200V264c0-110.5 89.5-200 200-200h496c110.5 0 200 89.5 200 200v496c0 110.5-89.5 200-200 200z',
				'path://M465 96.9L96.5 364.6c-28 20.4-39.8 56.5-29.1 89.4l140.7 433.1c10.7 33 41.4 55.3 76.1 55.3h455.4c34.7 0 65.4-22.3 76.1-55.3L956.5 454c10.7-33-1-69.1-29.1-89.4L559 96.9c-28-20.4-66-20.4-94 0z',
				'path://M675.9 107.2H348.1c-42.9 0-82.5 22.9-104 60.1L80 452.1c-21.4 37.1-21.4 82.7 0 119.8l164.1 284.8c21.4 37.2 61.1 60.1 104 60.1h327.8c42.9 0 82.5-22.9 104-60.1L944 571.9c21.4-37.1 21.4-82.7 0-119.8L779.9 167.3c-21.4-37.1-61.1-60.1-104-60.1z',
				'path://M1004.1 512L692 332 512 19.9 332 332 19.9 512 332 692l180 312.1L692 692z',
				'path://M541.1 102.4L662 347.3l270.2 39.3c26.7 3.9 37.3 36.6 18 55.4L754.7 632.5l46.2 269.1c4.6 26.5-23.3 46.8-47.1 34.3L512 808.8l-241.7 127c-23.8 12.5-51.7-7.7-47.1-34.3l46.2-269.1L73.8 442c-19.3-18.8-8.6-51.6 18-55.4L362 347.3l120.8-244.8c12-24.2 46.4-24.2 58.3-0.1z',
				'path://M796 512l104-180c16.9-29.3-4.2-65.9-38-65.9H654L550 86c-16.9-29.3-59.2-29.3-76.1 0L370 266.1H162.1c-33.8 0-54.9 36.6-38 65.9l104 180L124 692c-16.9 29.3 4.2 65.9 38 65.9h208L474 938c16.9 29.3 59.2 29.3 76.1 0l104-180.1H862c33.8 0 54.9-36.6 38-65.9L796 512z',
				'path://M818.7 477.5h-2.1c-16.8-120.7-120.4-213.7-245.8-213.7-89.9 0-168.6 47.8-212.2 119.4-28.1-15-60.2-23.6-94.4-23.6C153.7 359.6 64 449.3 64 559.9c0 107 84 194.5 189.7 200l565 0.3c78.1 0 141.3-63.3 141.3-141.3 0-78.1-63.3-141.4-141.3-141.4z',
				'path://M753.8 512l75-75c66.8-66.8 66.8-175.1 0-241.8-66.8-66.8-175.1-66.8-241.8 0l-75 75-75-75c-66.8-66.8-175.1-66.8-241.8 0-66.8 66.8-66.8 175.1 0 241.8l75 75-75 75c-66.8 66.8-66.8 175.1 0 241.8 66.8 66.8 175.1 66.8 241.8 0l75-75 75 75c66.8 66.8 175.1 66.8 241.8 0 66.8-66.8 66.8-175.1 0-241.8l-75-75z'
			]
			const selectAllIcon = 'path://M512 608a96 96 0 1 1 0-192 96 96 0 0 1 0 192m0-256c-88.224 0-160 71.776-160 160s71.776 160 160 160 160-71.776 160-160-71.776-160-160-160 M512 800c-212.064 0-384-256-384-288s171.936-288 384-288 384 256 384 288-171.936 288-384 288m0-640C265.248 160 64 443.008 64 512c0 68.992 201.248 352 448 352s448-283.008 448-352c0-68.992-201.248-352-448-352'
			const unSelectAllIcon = 'path://M512 800c-66.112 0-128.32-24.896-182.656-60.096l94.976-94.976A156.256 156.256 0 0 0 512 672c88.224 0 160-71.776 160-160a156.256 156.256 0 0 0-27.072-87.68l101.536-101.536C837.28 398.624 896 493.344 896 512c0 32-171.936 288-384 288m96-288a96 96 0 0 1-96 96c-14.784 0-28.64-3.616-41.088-9.664l127.424-127.424c6.048 12.448 9.664 26.304 9.664 41.088M128 512c0-32 171.936-288 384-288 66.112 0 128.32 24.896 182.656 60.096L277.536 701.216C186.72 625.376 128 530.656 128 512m664.064-234.816l91.328-91.328-45.248-45.248-97.632 97.632C673.472 192.704 595.456 160 512 160 265.248 160 64 443.008 64 512c0 39.392 65.728 148.416 167.936 234.816l-91.328 91.328 45.248 45.248 97.632-97.632C350.528 831.296 428.544 864 512 864c246.752 0 448-283.008 448-352 0-39.392-65.728-148.416-167.936-234.816 M512 352c-88.224 0-160 71.776-160 160 0 15.328 2.848 29.856 6.88 43.872l58.592-58.592a95.616 95.616 0 0 1 79.808-79.808l58.592-58.592A157.76 157.76 0 0 0 512 352'
	
			const LEGEND_LINE_HIGHT = 24
			const LEGEND_ITEM_FIXED = 8
			const LEGEND_CHAR_FIXED = 7.5
	
			function adjustedGridBottom(legendData) {
				var totalChars = legendData.reduce((total, item, a) => {
					return total + item.length
				}, 0)
				var legendWidth = $("#" + chartDomId).width()
				var lineCount = (totalChars * LEGEND_CHAR_FIXED + legendData.length * LEGEND_ITEM_FIXED) / legendWidth
				return Math.ceil(lineCount) * LEGEND_LINE_HIGHT
			}
	
			function getEchartOptionTemplate(id, title, legendData, labels, series, axisLabelRotate) {
				var legendGroupMap = {}
				var legendMap = {}
				legendData.forEach((name) => {
					var groupName = ''
					var parts = name.split('_')
					if (parts.length > 0) {
						groupName = parts[0]
					} else {
						groupName = name
					}
					if (!legendGroupMap.hasOwnProperty(groupName)) {
						legendGroupMap[groupName] = Object.keys(legendGroupMap).length % legendIcons.length
					}
					legendMap[name] = groupName
				})
				var newLegendData = legendData.map((name) => {
					return {
						name,
						icon: legendIcons[legendGroupMap[legendMap[name]]]
					}
				})
				var mySelectTools = {
					mySelectAll: {
						show: true,
						title: 'select all',
						icon: selectAllIcon,
						onclick: function () {
							var myChart = echarts.init(document.getElementById(id))
							var option = myChart.getOption()
							var selected = {}
							option.legend[0].data.forEach((item) => {
								selected[item.name] = true
							})
							option.legend[0].selected = selected
							myChart.setOption(option)
						}
					},
					myUnSelectAll: {
						show: true,
						title: 'unselect all',
						icon: unSelectAllIcon,
						onclick: function () {
							var myChart = echarts.init(document.getElementById(id))
							var option = myChart.getOption()
							var selected = {}
							option.legend[0].data.forEach((item) => {
								selected[item.name] = false
							})
							option.legend[0].selected = selected
							myChart.setOption(option)
						}
					}
				}
				for (var group in legendGroupMap) {
					mySelectTools['mySelectGroup' + group] = {
						show: true,
						title: group, //'select ' + group,
						icon: legendIcons[legendGroupMap[group]],
						onclick: function (parma, e, name) {
							var group = name.replace('mySelectGroup', '', 1)
							var myChart = echarts.init(document.getElementById(id))
							var option = myChart.getOption()
							var selected = {}
							option.legend[0].data.forEach((item) => {
								if (legendMap[item.name] == group) {
									selected[item.name] = true
								} else {
									selected[item.name] = false
								}
							})
							option.legend[0].selected = selected
							myChart.setOption(option)
						}
					}
				}
				var commonSelectTools = {
					magicType: {
						type: ['line', 'bar', 'stack', 'tiled']
					},
					dataZoom: {},
					restore: {},
					dataView: {},
					saveAsImage: {}
				}
				Object.assign(mySelectTools, commonSelectTools);
				return {
					title: {
						text: title,
						textStyle: {
							fontSize: 16
						}
					},
					tooltip: {
						trigger: 'axis',
						axisPointer: {
							type: 'cross',
							label: {
								backgroundColor: '#6a7985'
							}
						},
						textStyle: {
							fontSize: 12
						},
						position: function (point) {
							return point
						},
						formatter: function (params) {
							var content = params[0].axisValue + '<br/>'
							params.forEach(item => {
								var icon = '<svg viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" width="12" height="12"><path stroke="' + item.color + '" fill="' + item.color + '" d="' + legendIcons[legendGroupMap[legendMap[item.seriesName]]].replace('path://', '') + '"/></g></svg>'
								content += icon + ' ' + item.seriesName + ': ' + item.value + '<br/>'
							});
							return content;
						}
					},
					legend: {
						bottom: 0,
						data: newLegendData,
						selectedMode: 'multiple'
					},
					toolbox: {
						feature: mySelectTools
					},
					grid: {
						left: '3%',
						right: '4%',
						bottom: adjustedGridBottom(legendData),
						containLabel: true,
						borderColor: 'red'
					},
					xAxis: [
						{
							type: 'category',
							boundaryGap: false,
							data: labels,
							splitLine: {
								show: true,
								lineStyle: {
									color: "rgba(241, 238, 238, 1)"
								}
							},
							axisLabel: {
								rotate: axisLabelRotate
							}
						}
					],
					yAxis: [
						{
							type: 'value',
							splitLine: {
								show: true,
								lineStyle: {
									color: "rgba(241, 238, 238, 1)"
								}
							}
						}
					],
					series
				}
			}
			function getChartFromCSV(id, title, data) {
				var labels = []
				var series = []
				var legendData = []
				var table = ""
				var relArr = data.split("\n")
				if (!$.isEmptyObject(relArr) && relArr.length > 1) {
					for (var i = 0; i < relArr.length; i++) {
						var values = relArr[i];
						table += "<tr>"
						if (!$.isEmptyObject(values.trim())) {
							var objArr = values.trim().split(",");
							for (var j = 0; j < objArr.length; j++) {
								if (i == 0) {
									if (j == 0) {
										legendData = objArr
									} else {
										series[j - 1] = {
											name: objArr[j],
											data: [],
											type: 'line',
											symbolSize: 2
										}
									}
								} else {
									if (j == 0) {
										labels[i - 1] = objArr[j]
									} else {
										series[j - 1].data[i - 1] = objArr[j]
									}
								}
								if (j == 0) {
									if (i == 0) {
										table += "<th>index</th>"
									} else {
										table += "<th>" + i + "</th>"
									}
								}
								if (i == 0) {
									table += "<th>" + objArr[j] + "</th>"
								} else {
									table += "<td>" + objArr[j] + "</td>"
								}
							}
						}
						table += "</tr>"
					}
				}
				legendData.shift()
				var config = getEchartOptionTemplate(id, title, legendData, labels, series, 45)
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
							chartOption = getChartFromCSV(chartDomId, file.name, this.result)
							chart.clear()
							chart.setOption(chartOption.config)
							$("#canvas-perf-detail-table").html("<tbody>" + chartOption.table + "</tbody>")
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
	
			#canvas-perf-detail {
				width: 90%;
				height: 500px;
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
		<div class="perf-container">
			<div id="canvas-perf-detail"></div>
			<table id="canvas-perf-detail-table"></table>
		</div>
		<script>
			var csvResult = "{{.Data}}"
			var chartDomId = "canvas-perf-detail"
			var chartOption = getChartFromCSV(chartDomId, "", csvResult)
			var chart = echarts.init(document.getElementById(chartDomId), "light")
			chart.setOption(chartOption.config)
			$("#canvas-perf-detail-table").html("<tbody>" + chartOption.table + "</tbody>")
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
