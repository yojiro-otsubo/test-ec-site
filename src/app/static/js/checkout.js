let element = document.getElementById('payment-form');
const ClientSecret = element.dataset.secret;
const stripe = Stripe(element.dataset.pk);

const options = {
  clientSecret: ClientSecret,
  // Fully customizable with appearance API.
  appearance: {
    theme: 'stripe'
  },
};
  
// Set up Stripe.js and Elements to use in checkout form, passing the client secret obtained in step 2
const elements = stripe.elements(options);
  
// Create and mount the Payment Element
const paymentElement = elements.create('payment');
paymentElement.mount('#payment-element');



const form = document.getElementById('payment-form');

form.addEventListener('submit', async (event) => {
  event.preventDefault();

  const {error} = await stripe.confirmPayment({
    //`Elements` instance that was used to create the Payment Element
    elements,
    confirmParams: {
      return_url: 'https://my-site.com/order/123/complete',
    },
  });

  if (error) {
    // This point will only be reached if there is an immediate error when
    // confirming the payment. Show error to your customer (for example, payment
    // details incomplete)
    const messageContainer = document.querySelector('#error-message');
    messageContainer.textContent = error.message;
  } else {
    // Your customer will be redirected to your `return_url`. For some payment
    // methods like iDEAL, your customer will be redirected to an intermediate
    // site first to authorize the payment, then redirected to the `return_url`.
  }
});

