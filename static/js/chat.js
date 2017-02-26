$(".watch-edit").change(function () {
    $(this).attr("data-changed", "true");
});

$("#save").click(function(){
    var settings = {
        key_words: {},
        alert_times: {},
        new_key_words: [],
        new_alert_times: [],
        quotes_doc: parseInt($("#quotesdoc").val()),
    };
    if (!settings.quotes_doc) {
        settings.quotes_doc = 0
    }
    $("input[setting]").each(function() {
        settings[$(this).attr("setting")] = this.checked
    });
    $("#keyword_table tr:gt(0)").each(function() {
        var inputs = $(this).find('input');
        var keyIn = $(inputs[0]).val();
        var msgIn = $(inputs[1]).val();
        var idIn = $(inputs[0]).attr("data-id");
        var edited = $(inputs[0]).attr("data-changed") == "true" || $(inputs[1]).attr("data-changed") == "true";

        var keyword = {
            key: keyIn,
            msg: msgIn,
            id: parseInt(idIn)
        };

        if (edited && keyIn && msgIn && keyword.id) {
            settings.key_words[keyword.id] = false;
            keyword.id = 0;
            settings.new_key_words.push(keyword);
            return
        }

        if (keyword.id) {
            if (!keyIn || !msgIn) {
                settings.key_words[keyword.id] = false
            } else {
                settings.key_words[keyword.id] = true
            }
        } else {
            keyword.id = 0;
            settings.new_key_words.push(keyword);
        }

    });
    $("#alerttimes_table tr:gt(0)").each(function() {
            var inputs = $(this).find('input');
        var timeIn = $(inputs[0]).val();
        var msgIn = $(inputs[1]).val();
        var idIn = $(inputs[0]).attr("data-id");
        var edited = $(inputs[0]).attr("data-changed") == "true" || $(inputs[1]).attr("data-changed") == "true";


        var alerttime = {
            time: timeIn,
            msg: msgIn,
            id: idIn
        };

        if (edited && timeIn && msgIn && alerttime.id) {
            settings.alert_times[alerttime.id] = false;
            alerttime.id = 0;
            settings.new_alert_times.push(alerttime);
            return
        }

        if (alerttime.id) {
            if (!timeIn || !msgIn) {
                settings.alert_times[alerttime.id] = false
            } else {
                settings.alert_times[alerttime.id] = true
            }
        } else {
            alerttime.id = 0;
            settings.new_alert_times.push(alerttime);
        }
    });
    $.ajax({
        type: "PUT",
        url: $("#ChatId").text(),
        data: JSON.stringify(settings)
    }).success(function(data) {
        notification(1, "Success", "Saved Chat Settings");
    }).error(function(data, status) {
        notification(4, "Error", "");
    });
});

$("#add_keyword").click(function(){
    $('#keyword_table tr:last').after('<tr><td><input class="form-control" type="text" value=""></td><td><input class="form-control" type="text" value=""></td></tr>');
});

$("#add_alerttimes").click(function(){
    $('#alerttimes_table tr:last').after('<tr><td><input class="form-control" type="text" value="" placeholder="3:04PM -07"></td><td><input class="form-control" type="text" value=""></td></tr>');
});

$("#add_alerttime_now").click(function(){
    var nowTime = new Date();
    var getHour = nowTime.getHours();
    var getMinute = nowTime.getMinutes();

    var hour = getHour > 12 ? getHour - 12 : getHour;
    var halfClock = getHour > 12 ? 'PM' : 'AM';
    var minute = getMinute < 10 ? '0' + getMinute : getMinute;
    var tz = nowTime.getTimezoneOffset() / 60 * -1;
    tz = tz < 10 && tz >= 0 ? '0' + tz : tz;
    tz = tz > -10 && tz <= 0 ? '-0' + (tz * -1) : tz;

    $('#alerttimes_table tr:last').after('<tr><td><input class="form-control" type="text" value="' + hour + ':' + minute + halfClock + ' ' + tz + '" placeholder="3:04PM -07"></td><td><input class="form-control" type="text" value="NOW"></td></tr>');
});
