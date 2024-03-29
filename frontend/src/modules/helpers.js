// import ReactGA from "react-ga";

export const goToForm = () => {
  const anchor = document.querySelector('#contact-form-anchor');

  if (anchor) {
    anchor.scrollIntoView({
      behavior: 'smooth',
      block: 'center',
    });
  }
}

export const useAnalyticsEventTracker = (category="Blog category") => {
  // const eventTracker = (action = "test action", label = "test label") => {
  //   ReactGA.event({category, action, label});
  // }
  // return eventTracker;
}
