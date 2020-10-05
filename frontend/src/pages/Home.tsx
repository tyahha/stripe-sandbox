import React from "react";

import { loadStripe } from "@stripe/stripe-js";
const stripePromise = loadStripe(
  process.env.REACT_APP_STRIPE_PUBLIC_KEY ?? ""
);

export const Home = () => {
  const handleClick = async () => {
    const stripe = await stripePromise;

    const response = await fetch("/create-checkout-session", {method: "POST"})
    const session = await response.json();

    const result = await stripe?.redirectToCheckout({
      sessionId: session.id
    })

    if (result && result.error) {
      alert(result.error)
      console.error(result.error)
    }
  }

  return (
    <div>
      <button role="link" onClick={handleClick}>Checkout</button>
    </div>
  );
}
