window.onload = function () {
    const separatorFull = /;\s*/;
    const separatorPart = /,\s*/;
    const pieChart = $("#pieChart")
    const graphChart = $("#graphChart")
    const columnChart = $("#columnChart")

    let data = pieChart.attr('data-object');
    let datapoints = []

    data = data.substring(0, data.length - 1).split(separatorFull);
    for (let i = 0; i < data.length; i++) {
        datapoints[i] = {label: data[i].split(separatorPart)[0], y: data[i].split(separatorPart)[1]}
    }
    let optionsPie = {
        width: $( document ).width(),
        title: {
            text: "Полная стоимость проданных товаров по категориям"
        },
        subtitles: [{
            text: "За всё время"
        }],
        theme: "light2",
        exportEnabled: true,
        animationEnabled: true,
        data: [{
            type: "pie",
            startAngle: 40,
            toolTipContent: "<b>{label}</b>: {y}",
            showInLegend: "true",
            legendText: "{label}",
            indexLabelFontSize: 16,
            indexLabel: "{label} - {y} руб",
            dataPoints: datapoints
        }]
    };
    pieChart.CanvasJSChart(optionsPie);

    datapoints = [];
    data = graphChart.attr('data-object');
    data = data.substring(0, data.length - 1).split(separatorFull);
    let labels = [];
    for (let i = 0; i < data.length; i++) {
        datapoints[i] = {label: data[i].split(separatorPart)[0], y: data[i].split(separatorPart)[1], x: new Date(2024, parseInt(data[i].split(separatorPart)[2]))};
        if (!labels.includes(datapoints[i].label)) labels.push(datapoints[i].label);
    }
    let finalData = [];
    for (let i = 0; i < labels.length; i++) {
        let buffer = [];
        for (let j = 0; j < datapoints.length; j++) {
            if (labels[i] === datapoints[j].label) buffer.push({y: parseInt(datapoints[j].y), x: datapoints[j].x});
        }
        finalData[i] = {
            type: "splineArea",
            name: labels[i],
            showInLegend: true,
            yValueFormatString: "#,##0 руб",
            xValueFormatString: "MMM YYYY",
            dataPoints: buffer
        }
    }
    var optionsGraph = {
        width: $( document ).width()-50,
        exportEnabled: true,
        animationEnabled: true,
        title:{
            text: "Статистика продаж по категориям"
        },
        subtitles: [{
            text: "За 2024 год"
        }],
        theme: "light2",
        axisX: {
            title: 'Период',
            valueFormatString: "MMM YYYY"
        },
        axisY :{
            title: 'Продажи',
            includeZero: false,
            lineThickness: 0
        },
        toolTip: {
            shared: true
        },
        legend: {
            cursor: "pointer",
            itemclick: toggleDataSeries
        },
        data: finalData
    };
    graphChart.CanvasJSChart(optionsGraph);
    function toggleDataSeries(e) {
        if (typeof (e.dataSeries.visible) === "undefined" || e.dataSeries.visible) {
            e.dataSeries.visible = false;
        } else {
            e.dataSeries.visible = true;
        }
        e.chart.render();
    }

    datapoints = [];
    data = columnChart.attr('data-object');
    data = data.substring(0, data.length - 1).split(separatorFull);
    for (let i = 0; i < data.length; i++) {
        datapoints[i] = {label: data[i].split(separatorPart)[0], y: parseInt(data[i].split(separatorPart)[1])}
    }
    var optionsColumn = {
        width: $( document ).width(),
        exportEnabled: true,
        theme: "light2",
        title: {
            text: "Количество товаров на складе"
        },
        subtitles: [{
            text: "На текущий момент"
        }],
        axisY: {
            title: "Количество",
            suffix: " шт"
        },
        axisX: {
            title: "Товары"
        },
        data: [
            {
                // Change type to "doughnut", "line", "splineArea", etc.
                type: "column",
                dataPoints: datapoints
            }
        ]
    };

    $("#columnChart").CanvasJSChart(optionsColumn);
}