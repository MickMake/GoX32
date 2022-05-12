$(document).ready(function() {
    setTimeout(function() {
        $(".alert").alert('close');
    }, 2000);

    // $("#gain").slider({
    //     tooltip: 'always'
    // });
    //
    // function refresh_trigger() {
    //     $("#Trigger").load('/trigger');
    // }
    // setInterval(function() { refresh_trigger(); }, 1000);
    //
    // $('#select-all').click(function(event) {
    //     if(this.checked) {
    //         // Iterate each checkbox
    //         $(':checkbox').each(function() {
    //             this.checked = true;
    //         });
    //     } else {
    //         $(':checkbox').each(function() {
    //             this.checked = false;
    //         });
    //     }
    // });
});

function toggle(source) {
    var checkboxes = document.querySelectorAll('input[id="checkboxes"]');
    for (var i = 0; i < checkboxes.length; i++) {
        if (checkboxes[i] != source)
            checkboxes[i].checked = source.checked;
    }
}
