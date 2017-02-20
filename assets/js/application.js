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
});

function add_another_race() {
    // clone racesEntry div and then add right after it
    // used on event>new, event>edit and event>replicate
    $('#races').append($(".raceDetailForm:first").clone());
}

function checkDateRange() {
  start = document.getElementById("StartDate").value;
  // end = document.getElementById("EndDate").value;
  startD = document.getElementById("StartDate").valueAsDate;
  endD = document.getElementById("EndDate").valueAsDate;
  if (startD > endD) {
    document.getElementById("EndDate").value = start;
  }
}

function checkRegRange() {
  start = document.getElementById("RegOpenDate").value;
  // end = document.getElementById("EndDate").value;
  startD = document.getElementById("RegOpenDate").valueAsDate;
  endD = document.getElementById("RegCloseDate").valueAsDate;
  if (startD > endD) {
    document.getElementById("RegCloseDate").value = start;
  }
}
