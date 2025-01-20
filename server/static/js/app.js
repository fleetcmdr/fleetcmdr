
var AgentData = {};
AgentData.ID = 0;

function roundMe(n, sig) {
    if (n === 0) return 0;
    var mult = Math.pow(10, sig - Math.floor(Math.log(n < 0 ? -n: n) / Math.LN10) - 1);
    return Math.round(n * mult) / mult;
 }

function updateIndividualCPUProgressBars(thing) {
    // console.log(thing)

    $.each(thing.Extra, function(i,cluster){
        $.each(cluster.CPUs, function(j, cpu){
            divID = 'cpu_'+cluster.Name+'_'+cpu.CPU+'_activity'
            // console.log(divID)
            $('#'+divID+' .progress-bar').animate({width: (1-cpu.IdleRatio) * $('#'+divID).width()}, 500);

            freq = roundMe(cpu.FreqHz/1000000000, 3)
            $('#'+divID+' .progress-bar').html(freq + "Ghz");

            // if(thing.criticality == "nominal"){
            //     $('#'+divID+' .progress-bar').removeClass('bg-warning').removeClass('bg-critical').addClass('bg-success')
            // }
            // if(thing.criticality == "warning"){
            //     $('#'+divID+' .progress-bar').removeClass('bg-success').removeClass('bg-critical').addClass('bg-warning')
            // }
            // if(thing.criticality == "critical"){
            //     $('#'+divID+' .progress-bar').removeClass('bg-success').removeClass('bg-warning').addClass('bg-critical')
            // }
        })
    })
}

function updateProgressBar(divID, thing) {
    // console.log(divID)
    // console.log(thing)
    // console.log(thing.value)
    // divID = 'battery_percent'
    $('#'+divID+' .progress-bar').animate({width: (thing.Value/100) * $('#'+divID).width()}, 500);
    if(thing.Text != ""){
        $('#'+divID+' .progress-bar').html(thing.Text)
    }


    if(thing.criticality == "nominal"){
        $('#'+divID+' .progress-bar').removeClass('bg-warning').removeClass('bg-critical').addClass('bg-success')
    }
    if(thing.criticality == "warning"){
        $('#'+divID+' .progress-bar').removeClass('bg-success').removeClass('bg-critical').addClass('bg-warning')
    }
    if(thing.criticality == "critical"){
        $('#'+divID+' .progress-bar').removeClass('bg-success').removeClass('bg-warning').addClass('bg-critical')
    }

    if(divID == "cpu_activity"){
        updateIndividualCPUProgressBars(thing)
    }
}

