// node modules
import React from 'react';
import { render } from 'react-dom';
// import { Provider } from 'mobx-react';
import { configure } from 'mobx';

configure({
  computedRequiresReaction: true,
  enforceActions: 'observed',
});

const App = () => <div>Emona Textile</div>;

render(
    <App />,
  //<Provider {...STORES}>
    // <ThemeProvider>
    //   <Router />
    // </ThemeProvider>
  //</Provider>,
  window.document.getElementById('root'),
);
