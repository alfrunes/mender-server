// Copyright 2020 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
import React from 'react';

import { undefineds } from '../../../../../tests/mockData';
import { render } from '../../../../../tests/setupTests';
import { ArtifactUploadConfirmation } from './ArtifactUpload';

describe('ArtifactUpload Component', () => {
  it('renders correctly', async () => {
    const { baseElement } = render(
      <ArtifactUploadConfirmation
        creation={{ file: { name: 'testFile', size: 1024 }, type: 'mender' }}
        onboardingState={{ complete: true }}
        updateCreation={vi.fn()}
      />
    );
    const view = baseElement.firstChild;
    expect(view).toMatchSnapshot();
    expect(view).toEqual(expect.not.stringMatching(undefineds));
  });
});
