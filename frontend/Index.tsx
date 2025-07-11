import React from 'react';
import ReactDOM from 'react-dom/client';
import './styles.scss';
import App from './App';

function Index() {

	return (
		<div>Index</div>
	)
}

ReactDOM.createRoot(document.querySelector("#index")!).render(
	<React.StrictMode>
		<App />
	</React.StrictMode>
)