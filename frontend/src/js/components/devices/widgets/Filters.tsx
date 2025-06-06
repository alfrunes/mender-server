// Copyright 2015 Northern.tech AS
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
import React, { useCallback, useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { Add as AddIcon } from '@mui/icons-material';
// material ui
import { Button, Chip, Collapse } from '@mui/material';

import EnterpriseNotification from '@northern.tech/common-ui/EnterpriseNotification';
import { InfoHintContainer } from '@northern.tech/common-ui/InfoHint';
import MenderTooltip from '@northern.tech/common-ui/MenderTooltip';
import storeActions from '@northern.tech/store/actions';
import { BENEFITS, DEVICE_FILTERING_OPTIONS, emptyFilter } from '@northern.tech/store/constants';
import {
  getDeviceFilters,
  getFilterAttributes,
  getGlobalSettings,
  getIsEnterprise,
  getSelectedGroupInfo,
  getTenantCapabilities,
  getUserCapabilities
} from '@northern.tech/store/selectors';
import { getDeviceAttributes, saveGlobalSettings, setDeviceListState } from '@northern.tech/store/thunks';
import { filtersFilter } from '@northern.tech/store/utils';
import { deepCompare, toggle } from '@northern.tech/utils/helpers';

import FilterItem from './FilterItem';

const { setDeviceFilters } = storeActions;

export const getFilterLabelByKey = (key, attributes) => {
  const attr = attributes.find(attr => attr.key === key);
  return attr?.value ?? key ?? '';
};

const MAX_PREVIOUS_FILTERS_COUNT = 3;

export const Filters = ({ className = '', onGroupClick, open }) => {
  const [reset, setReset] = useState(false);
  const [newFilter, setNewFilter] = useState(emptyFilter);

  const dispatch = useDispatch();
  const { hasFullFiltering, plan } = useSelector(getTenantCapabilities);
  const { canManageUsers } = useSelector(getUserCapabilities);
  const { groupFilters, selectedGroup } = useSelector(getSelectedGroupInfo);
  const attributes = useSelector(getFilterAttributes);
  const filters = useSelector(getDeviceFilters);
  const isEnterprise = useSelector(getIsEnterprise);
  const { previousFilters = [] } = useSelector(getGlobalSettings);

  useEffect(() => {
    if (open) {
      dispatch(getDeviceAttributes());
    }
  }, [dispatch, open]);

  const saveUpdatedFilter = useCallback(
    updatedFilter => {
      if (canManageUsers && !previousFilters.find(filter => deepCompare(filter, updatedFilter))) {
        const changedPreviousFilters = [...previousFilters, updatedFilter];
        dispatch(saveGlobalSettings({ previousFilters: changedPreviousFilters.slice(-1 * MAX_PREVIOUS_FILTERS_COUNT) }));
      }
    },
    [canManageUsers, dispatch, previousFilters]
  );

  const handleFilterChange = useCallback(
    filters => {
      const activeFilters = filters.filter(filtersFilter).filter(item => item.value !== '');
      dispatch(setDeviceFilters(activeFilters));
      dispatch(setDeviceListState({ selectedId: undefined, page: 1, shouldSelectDevices: true, forceRefresh: true, filterSelection: undefined }));
    },
    [dispatch]
  );

  // We want to preview the resulting list while user types / selects a filter before saving
  const applyPreviewFilter = useCallback(
    updatedFilter => {
      const activeFilters = [...filters, updatedFilter].filter(filtersFilter).filter(item => item.key && item.value !== '');
      dispatch(setDeviceListState({ selectedId: undefined, page: 1, shouldSelectDevices: true, forceRefresh: true, filterSelection: activeFilters }));
    },
    [dispatch, filters]
  );

  const updateFilter = useCallback(
    updatedFilter => {
      saveUpdatedFilter(updatedFilter);
      handleFilterChange([...filters, updatedFilter]);
      setReset(toggle);
    },
    [filters, handleFilterChange, saveUpdatedFilter]
  );

  const resetIdFilter = () => dispatch(setDeviceListState({ selectedId: undefined, setOnly: true }));

  const removeFilter = removedFilter => {
    if (removedFilter.key === 'id') {
      resetIdFilter();
    }
    let changedFilters = filters.filter(filter => !deepCompare(filter, removedFilter));
    handleFilterChange(changedFilters);
  };

  const clearFilters = () => {
    handleFilterChange([]);
    resetIdFilter();
    setReset(toggle);
  };

  const onAddClick = () => updateFilter(newFilter);

  const isFilterDefined = Object.values(newFilter).every(thing => !!thing);
  const currentFilters = filters.filter(filtersFilter);
  return (
    <Collapse in={open} timeout="auto" className={`${className} filter-wrapper`} unmountOnExit>
      <>
        <div className="flexbox">
          <div className="margin-right" style={{ marginTop: currentFilters.length ? 8 : 25 }}>
            Devices matching:
          </div>
          <div>
            <div className="filter-list">
              {currentFilters.map(item => (
                <Chip
                  className="margin-right-small"
                  key={`filter-${item.key}-${item.operator}-${item.value}`}
                  label={`${getFilterLabelByKey(item.key, attributes)} ${DEVICE_FILTERING_OPTIONS[item.operator].shortform} ${
                    item.operator !== DEVICE_FILTERING_OPTIONS.$exists.key && item.operator !== DEVICE_FILTERING_OPTIONS.$nexists.key
                      ? item.operator === DEVICE_FILTERING_OPTIONS.$regex.key
                        ? `${item.value}.*`
                        : item.value
                      : ''
                  }`}
                  onDelete={() => removeFilter(item)}
                />
              ))}
            </div>
            {(hasFullFiltering || !currentFilters.length) && (
              <>
                <FilterItem attributes={attributes} onChange={setNewFilter} onSelect={applyPreviewFilter} onSave={updateFilter} plan={plan} reset={reset} />
                {isFilterDefined && <Chip className="margin-bottom-small" icon={<AddIcon />} label="Add a rule" color="primary" onClick={onAddClick} />}
              </>
            )}
            <EnterpriseNotification id={BENEFITS.fullFiltering.id} />
          </div>
        </div>
        {!!filters.length && !groupFilters.length && (
          <div className="flexbox column margin-top-small margin-bottom-small" style={{ alignItems: 'flex-end' }}>
            <span className="link margin-small margin-top-none" onClick={clearFilters}>
              Clear filter
            </span>
          </div>
        )}
        {isEnterprise && !!filters.length && (
          <div>
            {selectedGroup ? (
              !!groupFilters.length && (
                <MenderTooltip
                  title="Saved changes will not change the target devices of any ongoing deployments to this group, but will take effect for new deployments"
                  arrow
                >
                  <Button variant="contained" color="secondary" onClick={onGroupClick}>
                    Save group
                  </Button>
                </MenderTooltip>
              )
            ) : (
              <Button variant="contained" color="secondary" onClick={onGroupClick}>
                Create group with this filter
              </Button>
            )}
            <InfoHintContainer>
              <EnterpriseNotification id={BENEFITS.dynamicGroups.id} />
            </InfoHintContainer>
          </div>
        )}
      </>
    </Collapse>
  );
};

export default Filters;
