import { Code } from '@northern.tech/common-ui/copy-code';

import DeviceConfiguration from './eventdetails/deviceconfiguration';
import FileTransfer from './eventdetails/filetransfer';
import PortForward from './eventdetails/portforward';
import TerminalSession from './eventdetails/terminalsession';
import { UserChange } from './eventdetails/userchange';

const FallbackComponent = ({ item }) => {
  let content = '';
  try {
    content = JSON.stringify(item, null, 2);
  } catch (error) {
    content = `error parsing the logged event:\n${error}`;
  }
  return <Code style={{ whiteSpace: 'pre' }}>{content}</Code>;
};

const changeTypes = {
  user: 'user',
  device: 'device',
  tenant: 'tenant'
};

const configChangeDescriptor = {
  set_configuration: 'definition',
  deploy_configuration: 'deployment'
};

const EventDetailsDrawerContentMap = item => {
  const { type } = item.object || {};
  let content = { title: 'Entry details', content: FallbackComponent };
  if (type === changeTypes.user) {
    content = { title: `${item.action}d user`, content: UserChange };
  } else if (type === changeTypes.device && item.action.includes('terminal')) {
    content = { title: 'Remote session log', content: TerminalSession };
  } else if (type === changeTypes.device && item.action.includes('file')) {
    content = { title: 'File transfer', content: FileTransfer };
  } else if (type === changeTypes.device && item.action.includes('portforward')) {
    content = { title: 'Port forward', content: PortForward };
  } else if (type === changeTypes.device && item.action.includes('configuration')) {
    content = { title: `Device configuration ${configChangeDescriptor[item.action] || ''}`, content: DeviceConfiguration };
  } else if (type === changeTypes.device) {
    content = { title: 'Device change', content: FallbackComponent };
  } else if (type === changeTypes.tenant) {
    content = { title: `${item.action}d tenant`, content: UserChange };
  }
  return content;
};

export default EventDetailsDrawerContentMap;
