$("#save").click(function(){
    var settings = {
        key_words: [],
        alert_times: [],
    };
    $("input[setting]").each(function() {
        settings[$(this).attr("setting")] = this.checked
    });
    $("#keyword_table tr:gt(0)").each(function() {
        var inputs = $(this).find('input');
        var keyIn = $(inputs[0]).val()
        var msgIn = $(inputs[1]).val()
        if (!keyIn || !msgIn) {
            return;
        }
        var keyword = {
            key: keyIn,
            msg: msgIn,
        }
        settings.key_words.push(keyword);
    });
    $("#alerttimes_table tr:gt(0)").each(function() {
            var inputs = $(this).find('input');
            var timeIn = $(inputs[0]).val()
            var msgIn = $(inputs[1]).val()
            if (!timeIn || !msgIn) {
                return;
            }
            var keyword = {
                time: timeIn,
                msg: msgIn,
            }
            settings.alert_times.push(keyword);
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
    $('#alerttimes_table tr:last').after('<tr><td><input class="form-control" type="text" value="" placeholder="3:04PM MST"></td><td><input class="form-control" type="text" value=""></td></tr>');
});

$("#add_alerttime_now").click(function(){
    var nowTime = new Date()
    var getHour = nowTime.getHours()
    var getMinute = nowTime.getMinutes()

    var hour = getHour > 12 ? getHour - 12 : getHour;
    var halfClock = getHour > 12 ? 'PM' : 'AM';
    var minute = getMinute < 10 ? '0' + getMinute : getMinute;
    var tz = /\((.*)\)/.exec(nowTime.toString())[1];

    $('#alerttimes_table tr:last').after('<tr><td><input class="form-control" type="text" value="' + hour + ':' + minute + halfClock + ' ' + tz + '" placeholder="3:04PM MST"></td><td><input class="form-control" type="text" value="NOW"></td></tr>');
});
