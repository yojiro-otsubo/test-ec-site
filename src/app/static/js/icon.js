$('#icon-img-input').on('change', function (e) {
  var reader = new FileReader();
  reader.onload = function (e) {
      $("#icon-preview").attr('src', e.target.result);
  }
  reader.readAsDataURL(e.target.files[0]);
  $("#icon-preview").css('display', 'block');
});


$('.my-icon').on('click',function(){
  $('#icon-img-input').trigger("click");
});