
$('.savecontact').click(function(){

    var customerName=$.trim($("#customername").val());
    var customerEmail=$.trim($("#customeremail").val());

    if(customerName==''){
    alert('Please enter full name');
    $( "#customername" ).focus();
    return false;
    }else if(customerEmail==''){
    alert('Please enter email address');
    $( "#customeremail" ).focus();
    return false;
    }

    $(".contactcolor").removeClass('bg-primary-600');
    $(".contactcolor").addClass('bg-success-600');
    $(".paymentcolor").removeClass('bg-light-600');
    $(".paymentcolor").addClass('bg-primary-600');

    $(".contact-loader").html('<iconify-icon icon="svg-spinners:wind-toy" width="20" height="20"  style="color: #FFFFFF"></iconify-icon>');
    $( ".contactform" ).hide();
    $( ".contactdata" ).show();
    $( ".editbtn" ).show();
    $(".contact-loader").html('');
    $(".cname").html(customerName);
    $(".cemail").html(customerEmail);

});

$('.editbtn').click(function(){
$( ".editbtn" ).hide();
$( ".contactform" ).show();
$( ".contactdata" ).hide();
});
$(".process").click(function(){

    var sender_name = $.trim($(".cname").html());
    var sender_email = $.trim($(".cemail").html());
    
    // Set session storage item
    sessionStorage.setItem('sender_name', sender_name);
    sessionStorage.setItem('sender_email', sender_email);
    
        if(sender_name==''){
        alert('Please enter full name and click to continue');
        $( "#customername" ).focus();
        return false;
        }else if(sender_email==''){
        alert('Please enter email address and click to continue');
        $( "#customeremail" ).focus();
        return false;
        }
    
    
    var cid=$(this).attr('data-cid');
    var cryptoID=$(this).attr('data-tid');
    var price_currency = $("#Requestedcurrency").html().toLowerCase();
    var price_amount = $("#Requestedamount").html();
    var order_description = $("#Description").html();
    var Client_id = $("#CrypTox").attr('data-cid');
    var CommonURL = $("#CrypTox").attr('href');
    var iid = getUrlParameter('iid');
    
    if(cryptoID=="" || price_currency=="" || price_amount=="" || Client_id=="" || CommonURL=="" || cid==""|| iid==""){
        location.reload(true);
    }
    var data="";
    var formData = {
        cid: cid,
        price_currency: price_currency,
        price_amount: price_amount,
        sender_name: sender_name,
        sender_email: sender_email,
        client_id: Client_id,
        pay_type: 1,
        crypto_id: cryptoID,
        customerrefid: iid
    };
    //loader coinLoader
    $(".coinLoader-"+cryptoID).html('<iconify-icon icon="svg-spinners:wind-toy" width="15" height="15"  style="color: #16A34A"></iconify-icon>');
   // alert(JSON.stringify(formData));
    $.ajax({
    url: "/pay-data",
    data: $.param(formData),
    type: "POST",
    contentType: 'application/x-www-form-urlencoded',
    success:function(data){
    
    
    
    if(data.amount){
    
    $("#second").hide();$("#third").show();
    
            $("#getamountcoin").html(data.amount); 
            $("#getamountcoin1").html(data.amount);
            $("#getamountcurrency").html(price_amount + " " + price_currency);
            $("#getaddress").html(data.address);
            $(".coincode").html(cid);
            $("#coinicon").attr("src","/views/images/" + data.coinicon);
            $("#getnetworkid").html(data.coinnetwork);
            $("#transid").html(data.transid);
            ////////////////Mobile Intent/////////////////////////////
            const deviceType = detectDevice();
	        //alert(deviceType);
           // Payment parameters
           // alert(data.address);
           // alert(data.amount);
           // alert(cid);
            if(cid=="btc"){
             cid="bitcoin"
            }

            if(deviceType !=  "Web"){
                
   
           // Universal payment deep link (example)
           const paymentLink = `${cid}://${data.address}?amount=${encodeURIComponent(data.amount)}`;
          // alert(paymentLink);
   
           // Fallback store link if the app is not installed
           if(deviceType ==  "Android"){
           const fallbackLink = "https://play.google.com/store/search?q=crypto-app&hl=en"; // Replace with the wallet download page
           }else{
           const fallbackLink = "https://apps.apple.com/app/trust-crypto-bitcoin-wallet/id1288339409"; // Replace with the wallet download page
           }
           // Attempt to open the wallet app
           window.location.href = paymentLink;
   
           // Fallback to the wallet download page after a delay
           setTimeout(function () {
             window.location.href = fallbackLink;
           }, 5000); // 2-second delay before redirecting to the fallback
         }
    
            ///////////////////////////////////////////
           // var timeLeft = 600;
            var timeLeft = 100;
            var elem = document.getElementById('timer');
            var timerId = setInterval(countdown, 1000);
    
    function countdown() {
        if (timeLeft == -1) {
            clearTimeout(timerId);
            timeLeft--;
         } else if ((timeLeft == 75) || (timeLeft == 50) || (timeLeft == 25)|| (timeLeft == 5)) {
    
        //} else if ( (timeLeft == 570) || (timeLeft == 480) || (timeLeft == 420)|| (timeLeft == 360)|| (timeLeft == 300)|| (timeLeft == 240) || (timeLeft == 180) || (timeLeft == 120) || (timeLeft == 60)) {
          ////////////////
         
            var status_coin = $('.coincode').html();
            var status_address = $('#getaddress').html();
            var status_transid = $('#transid').html();
            
          var formDataStatus = {
            status_coin: status_coin,
            status_address: status_address,
            status_transid: status_transid,
            status_coinid: data.coin_id,
            client_id: Client_id
        };
        //alert(JSON.stringify(formDataStatus));
        
            $.ajax({
                url: "/check-payment-status",
                data: $.param(formDataStatus),
                type: "POST",
                contentType: 'application/x-www-form-urlencoded',
                success:function(status){
                    //alert(status.payment_status);
                    if(status.payment_status=="Success"){
                        url = CommonURL+"/success/"+status.payment_id;
                        $( location ).attr("href", url);
                    }else if(status.payment_status=="Declined"){
                            url = CommonURL+"/failed/"+status.payment_id;
                            $( location ).attr("href", url);}
                    else if(status.payment_status=="Dispute"){
                            url = CommonURL+"/dispute/"+status.payment_id;
                            $( location ).attr("href", url);
                    }else{
                        //alert(status.payment_status)
                    }
    
                },
                error:function (){}
                });	
    
        timeLeft--;
        
          ////////////////
        
        } else if (timeLeft == 1) {   
           // alert("UP")
            var status_transid = $('#transid').html();
            url =  CommonURL+"/failed/"+status_transid;
            $( location ).attr("href", url);
            timeLeft--;
        } else {
            var timeLeftNew =formatTime(timeLeft)
           // alert(vvv)
            elem.innerHTML = timeLeftNew ;
            //elem.innerHTML = timeLeft ;
            timeLeft--;
        }
    }
            //////////////////////////////////////////
    
            setTimeout(
                function() 
                {
                  //alert("Display Button")
                 // $("#statusbtn").show();
                  //$("#replybtn").hide();
                }, 10000);
    
                //$("#successlink").attr("href","{{$.CommonURL}}/success/"+data.transid); 
               // $("#failedlink").attr("href","{{$.CommonURL}}/success/"+data.transid);
    
    }else{
     //alert(99)
    // $(".returnfirst").click();
     }
    
    // generate QR code
    $(".qrcodes").qrcode({
        text: data.qr_code ,
        size: 150
    });
    
    
    },
    error:function (){}
    });	
    
    });
    
    $('#replyButton').click(function(){
      location.reload(true);
     // $(".process").off('click');
      //$("#second").show();
      //$("#third").hide();
      //$('.process').prop('disabled', true)
    });
    







// for hide footer
$(".d-footer").css("display", "none"); 