$(document).ready(function() {
  setTimeout(function() {
    $(".alert").alert('close');
  }, 2000);

  $('#record_sermon').on('click', function() {
    // alert($(this).prop('checked'));
    $('#config').submit();
  });
  $('#record_band').on('click', function() {
    // alert($(this).prop('checked'));
    $('#config').submit();
  });

});
