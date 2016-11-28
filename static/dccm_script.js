var values = [];
var labels = [];
var value = 0;
values.length = 100;
labels.length = 100;
values.fill(0);
labels.fill(0);

var canvas = document.getElementById("myChart"),
	ctx = canvas.getContext('2d'),
	indata = {
			labels: labels,
			datasets: [{
				label: '# of Votes',
				data: values,
				backgroundColor: 'rgba(0, 99, 255, 0.2)',
				borderColor: 'rgba(0,99,255,0.5)',
				borderWidth: 2,
				lineTension: 0.25,
				pointRadius: 0
        }]
    },	
	inoption = {
		responsive: false,
		animation: {
			duration: 250*1,
			easing: 'linear'
		},
		legend: false,
	    scales: {
		    yAxes: [{
                ticks: {
                    beginAtZero:true,
					max: 100,
					min: 0
                }
            }],
			 xAxes: [{
				display: false
				
            }]  			
        }
    };
	
var barChart = new Chart(ctx, {type: 'line', data: indata, options: inoption});

  
function chart_update() {
		
	var x = new XMLHttpRequest();	
	
		
	x.onreadystatechange = function() {
        if(this.readyState == 4 && this.status == 200) {
			
          value = this.responseText;
		   	  
        }        
	};
	
	x.open("GET", "/getCPUusage", true);
	x.send();
	
	values.push(value);
	values.shift();
	
	barChart.update();
 
 }

setInterval(function(){
	requestAnimationFrame(chart_update);  
}, 250);
