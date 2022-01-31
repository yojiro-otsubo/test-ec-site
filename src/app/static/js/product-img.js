$('.none-active1').on('click',function(){
    $('#p-img-1').css('display', 'block');
    $('#p-img-2').css('display', 'none');
    $('#p-img-3').css('display', 'none');

    $('#p-img-soldout-1').css('display', 'block');
    $('#p-img-soldout-2').css('display', 'none');
    $('#p-img-soldout-3').css('display', 'none');

    $('.none-active1').css('border', '1px solid rgb(0, 174, 255)');
    $('.none-active2').css('border', '0.5px solid gainsboro');
    $('.none-active3').css('border', '0.5px solid gainsboro');
  });


  $('.none-active2').on('click',function(){
    $('#p-img-1').css('display', 'none');
    $('#p-img-2').css('display', 'block');
    $('#p-img-3').css('display', 'none');

    $('#p-img-soldout-1').css('display', 'none');
    $('#p-img-soldout-2').css('display', 'block');
    $('#p-img-soldout-3').css('display', 'none');

    $('.none-active1').css('border', '0.5px solid gainsboro');
    $('.none-active2').css('border', '1px solid rgb(0, 174, 255)');
    $('.none-active3').css('border', '0.5px solid gainsboro');
  });

  $('.none-active3').on('click',function(){
    $('#p-img-1').css('display', 'none');
    $('#p-img-2').css('display', 'none');
    $('#p-img-3').css('display', 'block');

    $('#p-img-soldout-1').css('display', 'none');
    $('#p-img-soldout-2').css('display', 'none');
    $('#p-img-soldout-3').css('display', 'block');

    $('.none-active1').css('border', '0.5px solid gainsboro');
    $('.none-active2').css('border', '0.5px solid gainsboro');
    $('.none-active3').css('border', '1px solid rgb(0, 174, 255)');
  });