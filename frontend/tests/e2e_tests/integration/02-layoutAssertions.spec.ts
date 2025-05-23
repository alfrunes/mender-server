// Copyright 2021 Northern.tech AS
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
import test, { expect } from '../fixtures/fixtures.ts';
import { selectors, timeouts } from '../utils/constants.ts';

test.describe('Layout assertions', () => {
  let navbar;
  test.beforeEach(async ({ page }) => {
    navbar = page.locator('.leftFixed.leftNav');
  });
  test.describe('Overall layout and structure', () => {
    test('shows the left navigation', async () => {
      await expect(navbar.getByRole('link', { name: /Dashboard/i })).toBeVisible();
      await expect(navbar.getByRole('link', { name: /Devices/i })).toBeVisible();
      await expect(navbar.getByRole('link', { name: /Releases/i })).toBeVisible();
      await expect(navbar.getByRole('link', { name: /Deployments/i })).toBeVisible();
    });

    test('has clickable header buttons', async () => {
      await expect(navbar.getByRole('link', { name: /Dashboard/i })).toBeVisible();
      await navbar.getByRole('link', { name: /Dashboard/i }).click();
      await navbar.getByRole('link', { name: /Devices/i }).click();
      await navbar.getByRole('link', { name: /Releases/i }).click();
      await navbar.getByRole('link', { name: /Deployments/i }).click();
    });
  });

  test('can authorize a device', async ({ browserName, page }) => {
    // allow twice the device interaction time + roughly a regular test execution time
    test.setTimeout(2 * timeouts.sixtySeconds + timeouts.fifteenSeconds);
    await navbar.getByRole('link', { name: /Devices/i }).click({ force: browserName === 'webkit' });
    let hasAcceptedDevice = false;
    try {
      await page.waitForSelector(`css=${selectors.deviceListItem}`, { timeout: timeouts.default });
      hasAcceptedDevice = await page.isVisible(selectors.deviceListItem);
    } catch {
      console.log(`no accepted device present so far`);
    }
    if (!hasAcceptedDevice) {
      const pendingMessage = await page.getByText(/pending authorization/i);
      await pendingMessage.waitFor({ timeout: timeouts.sixtySeconds });
      await pendingMessage.click();
      await page.click(selectors.deviceListCheckbox);
      await page.click('.MuiSpeedDial-fab');
      await page.click('[aria-label="accept"]');
    }
    await page.locator(`input:near(:text("Status:"))`).first().click({ force: true });
    await page.click(`css=.MuiPaper-root >> text=/Accepted/i`);
    await page.waitForSelector(`css=${selectors.deviceListItem} >> text=/original/`, { timeout: timeouts.sixtySeconds });
    const element = await page.textContent(selectors.deviceListItem);
    expect(element.includes('original')).toBeTruthy();
    await page.locator(`css=${selectors.deviceListItem} div:last-child`).last().click();
    await page.getByText(/Device information for/i).waitFor();
    await expect(page.getByText('Authentication status')).toBeVisible();
  });

  test('can group a device', async ({ browserName, page }) => {
    await navbar.getByRole('link', { name: /Devices/i }).click({ force: browserName === 'webkit' });
    const groupList = await page.locator('.grouplist');
    const wasGrouped = await groupList.getByText('testgroup').isVisible();
    test.skip(wasGrouped, 'looks like the device was grouped already, continue with the remaining tests');
    await page.click(selectors.deviceListCheckbox);
    await page.click('.MuiSpeedDial-fab');
    await page.click('[aria-label="group-add"]');
    await page.getByLabel(/type to create new/i).fill('testgroup');
    await page.click('.MuiDialogTitle-root');
    const groupCreation = await page.getByRole('button', { name: /create group/i });
    const groupExtension = await page.getByRole('button', { name: /add to group/i });
    await groupCreation.or(groupExtension).first().click();
    await groupList.getByText('testgroup').waitFor();
    await expect(groupList.getByText('testgroup')).toBeVisible();
    await groupList.getByText('All devices');
    await page.click(selectors.deviceListCheckbox);
    await groupList.getByText('testgroup').click();
    await expect(page.locator(`css=${selectors.deviceListItem} >> text=/original/`)).toBeVisible();
  });
});
