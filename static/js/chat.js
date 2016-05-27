$("#save").click(function(){
    var settings = {};
    $("input[setting]").each(function() {
        settings[$(this).attr("setting")] = this.checked
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
