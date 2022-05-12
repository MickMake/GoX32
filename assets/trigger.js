$(document).ready(function() {
    setTimeout(function() {
        $(".alert").alert('close');
    }, 2000);

    $('#record_sermon').on('click', function() {
        alert($(this).prop('checked'));
        $('#config').submit();
    });

    $('#record_band').on('click', function() {
        alert($(this).prop('checked'));
        $('#config').submit();
    });

    function strcmp(a, b) {
        // return (a<b?-1:(a>b?1:0));
        return (a === b);
    }

    function refresh_trigger() {
        $("#LastArgs").load("/lastargs");
        args = $("#LastArgs").text();

        $("#LastCmd").load("/lastcmd");
        cmd = $("#LastCmd").text();

        if (strcmp(args, "ON")) {
            $("#ON").addClass("highlight");
            $("#OFF").removeClass("highlight");
            $('input[name="Trigger"]').val(cmd);
        } else if (strcmp(args, "OFF")) {
            $("#ON").removeClass("highlight");
            $("#OFF").addClass("highlight");
            $('input[name="Trigger"]').val(cmd);
        } else {
            $("#ON").removeClass("highlight");
            $("#OFF").removeClass("highlight");
            $('input[name="Trigger"]').val("");
        }
    }
    setInterval(function() { refresh_trigger(); }, 200);
});
