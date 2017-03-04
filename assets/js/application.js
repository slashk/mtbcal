require('expose?$!expose?jQuery!jquery');
require("bootstrap/dist/js/bootstrap.js");

$(() => {
    $('[data-toggle="tooltip"]').tooltip();

    // Make sure EndDate is after StartDate
    if ($('#StartDate').length) {
      $('#StartDate').blur(function() {
        checkDateRange();
    })
  }

    // Make sure RegCloseDate is after RegOpenDate
    if ($('#RegOpenDate').length) {
      $('#RegOpenDate').blur(function() {
        checkRegRange();
      });
    };
    $('#add_race_button').click(function () {
        add_another_race();
    });
    $('#remove_race_button').click(function () {
        remove_race();
    });
});

function add_another_race() {
  var template = $('#races .raceDetailForm:first').clone();
  raceCount = $('.raceDetailForm').length;
  var section = template.clone().find(':input').each(function(){
    var fieldName = $('label[for="' + this.name + '"]').html();
    var newId = "Race." + raceCount + "." + fieldName;
    // var newName = "Race." + raceCount + "." + fieldName;
    this.id = newId;
    this.name = newId;
  }).end()
  .appendTo('#races');
  return false;
}

function remove_race() {
  $('#races .raceDetailForm:last').fadeOut(300, function() { $(this).remove(); });
  return false;
}

function checkDateRange() {
  start = document.getElementById("StartDate").value;
  startD = document.getElementById("StartDate").valueAsDate;
  endD = document.getElementById("EndDate").valueAsDate;
  if (startD > endD) {
    document.getElementById("EndDate").value = start;
  }
}

function checkRegRange() {
  start = document.getElementById("RegOpenDate").value;
  startD = document.getElementById("RegOpenDate").valueAsDate;
  endD = document.getElementById("RegCloseDate").valueAsDate;
  if (startD > endD) {
    document.getElementById("RegCloseDate").value = start;
  }
}
