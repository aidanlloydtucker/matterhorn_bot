/**
 * Creates a notification.
 *
 * @constructor
 * @param {number} level The level of alert. 1=success, 2=info, 3=warning, 4=danger.
 * @param {string} title The title of the alert. Bolded.
 * @param {string} text The message of the alert.
 */
function notification(level, title, text) {
    switch (level) {
        case 1:
            level = "success";
            break;
        case 2:
            level = "info";
            break;
        case 3:
            level = "warning";
            break;
        case 4:
            level = "danger";
            break;
        default:
            level = "success";
    }
    var id = Math.floor((Math.random() * 100) + 1);
    $("body").prepend('<div id="notification' + id + '" class="alert alert-'+level+' animated"> <div class="container"> <b>' + title + '</b> ' + text + '</div>');
    $("#notification" + id).addClass('fadeInLeft');
    setTimeout(function() {
        $("#notification" + id).addClass('fadeOutRight');
        setTimeout(function() {
            $("#notification" + id).remove();
        },1000);
    },5000);
}
