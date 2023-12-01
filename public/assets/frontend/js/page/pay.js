$(function () {
    var timer, minutes, seconds, ci, qi;

    var timeDown = function($container){
        timer = parseInt(remainSeconds) - 1;
        ci = setInterval(function () {
            minutes = parseInt(timer / 60, 10)
            seconds = parseInt(timer % 60, 10);
            if(minutes<=0){
                minutes="00";
            }else{
                minutes = minutes < 10 ? "0" + minutes : minutes;
            }
            if(seconds<=0){
                seconds="00";
            }else{
                seconds = seconds < 10 ? "0" + seconds : seconds;
            }
            $container.find(".minutes b").text(minutes);
            $container.find(".seconds b").text(seconds);
            if (--timer < 0) {
                $("#qrcodebox .expired").removeClass("hidden");
                var msg = messages.expire;
                if(cancelURL) aAttr = '<a href="'+cancelURL+'" class="text-danger">' + messages.expire + '</a>';
                $("#msg-summary").removeClass("alert-warning").addClass("alert-danger").html(msg).removeClass("hidden");
                clearInterval(ci);
                clearTimeout(qi);
            }
        }, 1000);
    }

    //定时查询订单状态
    var checkOrderStatus = function () {
        if(qi) clearTimeout(qi);
        $.ajax({
            url: queryURL,
            dataType: 'json',
            success: function (ret) {
                if (ret.Code < 1) {
                    //$("#msg-summary").html(ret.Info).removeClass("hidden");
                    qi = setTimeout(checkOrderStatus, 3000);
                    return;
                }
                var data = ret.Data;
                switch(data.status){
                    case 'expired':
                        $("#qrcodebox .expired").removeClass("hidden");
                        return;
                    case 'inprogress':
                        break;
                    case 'success':
                        clearTimeout(ci);
                        var msg;
                        $("#qrcodebox .paid").removeClass("hidden");
                        if (data.returnURL != '') {
                            setTimeout(function () {
                                location.href = data.returnURL;
                            }, 2000);
                            msg=messages.jump;
                        } else {
                            msg=messages.success;
                        }
                        $("#msg-summary").removeClass("alert-warning").addClass("alert-success").html(msg).removeClass("hidden");
                        return;
                    default:
                        var msg = data.reason;
                        if(!msg) msg = 'Unknown Error';
                        $("#msg-summary").html(msg).removeClass("hidden");
                        return;
                }
                qi = setTimeout(checkOrderStatus, 3000);
            },
            error: function () {
                qi = setTimeout(checkOrderStatus, 3000);
            }
        });
    };
    if($('#remainseconds').length>0) timeDown($('#remainseconds'));
    if(queryURL) checkOrderStatus();
});