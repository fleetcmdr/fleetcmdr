
var AgentData = {};
AgentData.ID = 0;

function updateCPUChart(thing) {
    // console.log(thing)
    $('#cpu_activity .progress-bar').animate({width: (thing.cpu/100) * $('#cpu_activity').width()}, 1000);

    if(thing.cpu < 80){
        $('#cpu_activity .progress-bar').removeClass('bg-warning').removeClass('bg-critical').addClass('bg-success')
    }
    if(thing.cpu >= 80 && thing.cpu < 90){
        $('#cpu_activity .progress-bar').removeClass('bg-success').removeClass('bg-critical').addClass('bg-warning')
    }
    if(thing.cpu >= 90){
        $('#cpu_activity .progress-bar').removeClass('bg-success').removeClass('bg-warning').addClass('bg-critical')
    }
}

function updateRAMChart(thing) {
    // console.log(thing)
    $('#ram_pressure .progress-bar').animate({width: (thing.ram/100) * $('#ram_pressure').width()}, 1000);

    if(thing.ram < 80){
        $('#ram_pressure .progress-bar').removeClass('bg-warning').removeClass('bg-critical').addClass('bg-success')
    }
    if(thing.ram >= 80 && thing.ram < 90){
        $('#ram_pressure .progress-bar').removeClass('bg-success').removeClass('bg-critical').addClass('bg-warning')
    }
    if(thing.ram >= 90){
        $('#ram_pressure .progress-bar').removeClass('bg-success').removeClass('bg-warning').addClass('bg-critical')
    }
}

function updateDiskChart(thing) {
    // console.log(thing)

    $('#disk_consumed .progress-bar').animate({width: (thing.disk/100) * $('#disk_consumed').width()}, 250);

    if(thing.disk < 80){
        $('#disk_consumed .progress-bar').removeClass('bg-warning').removeClass('bg-critical').addClass('bg-success')
    }
    if(thing.disk >= 80 && thing.disk < 90){
        $('#disk_consumed .progress-bar').removeClass('bg-success').removeClass('bg-critical').addClass('bg-warning')
    }
    if(thing.disk >= 90){
        $('#disk_consumed .progress-bar').removeClass('bg-success').removeClass('bg-warning').addClass('bg-critical')
    }
}
